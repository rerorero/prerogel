package worker

import (
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/aggregator"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus"
)

type lastAggregated struct {
	superstep uint64
	values    map[string]*any.Any
}

type coordinatorActor struct {
	util.ActorUtil
	behavior              actor.Behavior
	plugin                plugin.Plugin
	workerProps           *actor.Props
	clusterInfo           *command.ClusterInfo
	ackRecorder           *util.AckRecorder
	aggregatedCurrentStep map[string]*any.Any
	lastAggregatedValue   lastAggregated
	currentStep           uint64
}

// NewCoordinatorActor returns an actor instance
func NewCoordinatorActor(plg plugin.Plugin, workerProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &util.AckRecorder{}
	ar.Clear()
	a := &coordinatorActor{
		plugin: plg,
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		workerProps: workerProps,
		ackRecorder: ar,
	}
	a.behavior.Become(a.idle)
	return a
}

// Receive is message handler
func (state *coordinatorActor) Receive(context actor.Context) {
	if state.ActorUtil.IsSystemMessage(context.Message()) {
		// ignore
		return
	}

	switch context.Message().(type) {
	case *command.CoordinatorStats:
		s := &command.CoordinatorStatsAck{}
		if state.lastAggregatedValue.values != nil {
			stats, err := state.getStats(state.lastAggregatedValue.values)
			if err != nil {
				state.ActorUtil.Fail(err)
				return
			}
			s.SuperStep = state.lastAggregatedValue.superstep
			s.NrOfActiveVertex = stats.ActiveVertices
			s.NrOfSentMessages = stats.MessagesSent
		}
		context.Respond(s)
		return

	default:
		state.behavior.Receive(context)
	}
}

func (state *coordinatorActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.NewCluster:
		if state.clusterInfo != nil {
			state.ActorUtil.Fail(fmt.Errorf("cluster info has already been set: %+v", state.clusterInfo))
			return
		}

		assigned, err := assignPartition(len(cmd.Workers), cmd.NrOfPartitions)
		if err != nil {
			state.ActorUtil.Fail(err)
			return
		}
		state.ackRecorder.Clear()

		ci := &command.ClusterInfo{
			WorkerInfo: make([]*command.ClusterInfo_WorkerInfo, len(cmd.Workers)),
		}
		for i, wreq := range cmd.Workers {
			var pid *actor.PID
			if wreq.Remote {
				// remote actor
				pidRes, err := remote.SpawnNamed(wreq.HostAndPort, RemoteCoordinatorName, RemoteWorkerName, 30*time.Second)
				if err != nil {
					state.ActorUtil.Fail(errors.Wrapf(err, "failed to spawn remote actor: code=%v", pidRes.StatusCode))
					return
				}
				pid = pidRes.Pid
			} else {
				// local actor
				pid = context.Spawn(state.workerProps)
			}

			context.Request(pid, &command.InitWorker{
				Partitions: assigned[i],
			})
			state.ackRecorder.AddToWaitList(pid.GetId())
			ci.WorkerInfo[i] = &command.ClusterInfo_WorkerInfo{
				WorkerPid:  pid,
				Partitions: assigned[i],
			}
		}

		state.clusterInfo = ci
		state.ActorUtil.Logger.Debug("start initializing workers")
		return

	case *command.InitWorkerAck:
		if ok := state.ackRecorder.Ack(cmd.WorkerPid.GetId()); !ok {
			state.ActorUtil.LogError(fmt.Sprintf("InitWorkerAck from unknown worker: %v", cmd.WorkerPid))
			return
		}
		if state.ackRecorder.HasCompleted() {
			state.ackRecorder.Clear()
			state.aggregatedCurrentStep = make(map[string]*any.Any)
			state.currentStep = 0
			for _, wi := range state.clusterInfo.WorkerInfo {
				context.Request(wi.WorkerPid, &command.SuperStepBarrier{
					ClusterInfo: state.clusterInfo,
				})
				state.ackRecorder.AddToWaitList(wi.WorkerPid.GetId())
			}
			// TODO: handle worker timeout
			state.behavior.Become(state.superstep)
			state.ActorUtil.Logger.Debug("superstep barrier")
		}
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[idle] unhandled corrdinator command: command=%#v", cmd))
		return
	}
}

func (state *coordinatorActor) superstep(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrierWorkerAck:
		if ok := state.ackRecorder.Ack(cmd.WorkerPid.GetId()); !ok {
			state.ActorUtil.LogError(fmt.Sprintf("superstep barrier ack from unknown worker: %v", cmd.WorkerPid))
			return
		}
		if state.ackRecorder.HasCompleted() {
			state.ackRecorder.Clear()
			for _, wi := range state.clusterInfo.WorkerInfo {
				context.Request(wi.WorkerPid, &command.Compute{
					SuperStep:        state.currentStep,
					AggregatedValues: state.lastAggregatedValue.values,
				})
				state.ackRecorder.AddToWaitList(wi.WorkerPid.GetId())
			}
			state.behavior.Become(state.computing)
			state.ActorUtil.Logger.Debug(fmt.Sprintf("start computing: step=%v", state.currentStep))
		}
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[superstep] unhandled corrdinator command: command=%#v", cmd))
		return
	}
}
func (state *coordinatorActor) computing(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.ComputeWorkerAck:
		if ok := state.ackRecorder.Ack(cmd.WorkerPid.GetId()); !ok {
			state.ActorUtil.LogError(fmt.Sprintf("compute ack from unknown worker: %v", cmd.WorkerPid))
			return
		}

		if cmd.AggregatedValues != nil {
			if err := aggregateValueMap(state.plugin.GetAggregators(), state.aggregatedCurrentStep, cmd.AggregatedValues); err != nil {
				state.ActorUtil.Fail(err)
				return
			}
		}
		if state.ackRecorder.HasCompleted() {
			state.ackRecorder.Clear()

			// check if there are active vertices
			stats, err := state.getStats(state.aggregatedCurrentStep)
			if err != nil {
				state.ActorUtil.Fail(err)
				return
			}

			if stats.ActiveVertices == 0 {
				// finish superstep
				state.behavior.Become(state.idle)
				state.ActorUtil.Logger.Info(fmt.Sprintf("finish computing: step=%v", state.currentStep))

			} else {
				// move step forward
				state.currentStep += uint64(1)
				for _, wi := range state.clusterInfo.WorkerInfo {
					context.Request(wi.WorkerPid, &command.SuperStepBarrier{
						ClusterInfo: state.clusterInfo,
					})
					state.ackRecorder.AddToWaitList(wi.WorkerPid.GetId())
				}
				// TODO: handle worker timeout
				state.behavior.Become(state.superstep)
				state.ActorUtil.Logger.Debug(fmt.Sprintf("start computing: step=%v", state.currentStep))
			}

			// update aggregated values
			state.lastAggregatedValue.superstep = state.currentStep
			state.lastAggregatedValue.values = state.aggregatedCurrentStep
			state.aggregatedCurrentStep = make(map[string]*any.Any)
		}
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[computing] unhandled corrdinator command: command=%#v", cmd))
		return
	}
}

func (state *coordinatorActor) getStats(aggregated map[string]*any.Any) (*aggregator.VertexStats, error) {
	v, err := getAggregatedValue(state.plugin.GetAggregators(), aggregated, aggregator.VertexStatsName)
	if err != nil {
		return nil, err
	}

	stats, ok := v.(*aggregator.VertexStats)
	if !ok {
		return nil, fmt.Errorf("not VertexStats %#v", v)
	}

	return stats, nil
}

func assignPartition(nrOfWorkers int, nrOfPartitions uint64) ([][]uint64, error) {
	if nrOfWorkers == 0 {
		return nil, errors.New("no available workers")
	}
	if nrOfPartitions == 0 {
		return nil, errors.New("no partitions")
	}
	max := nrOfPartitions / uint64(nrOfWorkers)
	surplus := nrOfPartitions % uint64(nrOfWorkers)
	pairs := make([][]uint64, nrOfWorkers)

	var part, to uint64
	for i := range pairs {
		to = part + max - 1
		if surplus > 0 {
			to++
			surplus--
		}
		if to > nrOfPartitions {
			to = nrOfPartitions - 1
		}
		parts := make([]uint64, (int)(to-part+1))
		for j := range parts {
			parts[j] = part + uint64(j)
		}

		pairs[i] = parts
		part += uint64(len(parts))
	}

	return pairs, nil
}
