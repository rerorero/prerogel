package worker

import (
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/aggregator"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus"
)

type lastAggregated struct {
	superstep uint64
	values    map[string]*types.Any
}

type coordinatorActor struct {
	util.ActorUtil
	behavior              actor.Behavior
	plugin                plugin.Plugin
	workerProps           *actor.Props
	clusterInfo           *command.ClusterInfo
	ackRecorder           *util.AckRecorder
	aggregatedCurrentStep map[string]*types.Any
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
	a.behavior.Become(a.setup)
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
				state.ActorUtil.Fail(context, err)
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

func (state *coordinatorActor) setup(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.NewCluster:
		if state.clusterInfo != nil {
			state.ActorUtil.Fail(context, fmt.Errorf("cluster info has already been set: %+v", state.clusterInfo))
			return
		}

		assigned, err := assignPartition(len(cmd.Workers), cmd.NrOfPartitions)
		if err != nil {
			state.ActorUtil.Fail(context, err)
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
				pidRes, err := remote.SpawnNamed(wreq.HostAndPort, fmt.Sprintf("worker-%d", i), WorkerActorKind, 30*time.Second)
				if err != nil {
					state.ActorUtil.Fail(context, errors.Wrapf(err, "failed to spawn remote actor: code=%v", pidRes.StatusCode))
					return
				}
				pid = pidRes.Pid
			} else {
				// local actor
				pid = context.Spawn(state.workerProps)
			}

			context.Request(pid, &command.InitWorker{
				Coordinator: context.Self(),
				Partitions:  assigned[i],
			})
			state.ackRecorder.AddToWaitList(pid.GetId())
			ci.WorkerInfo[i] = &command.ClusterInfo_WorkerInfo{
				WorkerPid:  pid,
				Partitions: assigned[i],
			}
		}

		state.clusterInfo = ci

		context.Respond(&command.NewClusterAck{})
		state.ActorUtil.LogDebug(context, "start initializing workers")
		return

	case *command.InitWorkerAck:
		if ok := state.ackRecorder.Ack(cmd.WorkerPid.GetId()); !ok {
			state.ActorUtil.LogError(context, fmt.Sprintf("InitWorkerAck from unknown worker: %v", cmd.WorkerPid))
			return
		}
		if state.ackRecorder.HasCompleted() {
			state.ackRecorder.Clear()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogDebug(context, "become idle")
		}
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[setup] unhandled corrdinator command: command=%#v", cmd))
		return
	}
}

func (state *coordinatorActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.LoadVertex:
		w := state.findWorkerInfoByVertex(context, plugin.VertexID(cmd.VertexId))
		if w == nil {
			state.ActorUtil.LogError(context, fmt.Sprintf("couldn't find worker to assign: vertex=%v", cmd.VertexId))
			return
		}
		context.Forward(w.WorkerPid)
		return

	case *command.StartSuperStep:
		state.aggregatedCurrentStep = make(map[string]*types.Any)
		state.currentStep = 0
		for _, wi := range state.clusterInfo.WorkerInfo {
			context.Request(wi.WorkerPid, &command.SuperStepBarrier{
				ClusterInfo: state.clusterInfo,
			})
			state.ackRecorder.AddToWaitList(wi.WorkerPid.GetId())
		}
		// TODO: handle worker timeout
		state.behavior.Become(state.superstep)
		state.ActorUtil.LogDebug(context, "become superstep")
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[setup] unhandled corrdinator command: command=%#v", cmd))
		return
	}
}

func (state *coordinatorActor) superstep(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrierWorkerAck:
		if ok := state.ackRecorder.Ack(cmd.WorkerPid.GetId()); !ok {
			state.ActorUtil.LogError(context, fmt.Sprintf("superstep barrier ack from unknown worker: %v", cmd.WorkerPid))
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
			state.ActorUtil.LogDebug(context, fmt.Sprintf("start computing: step=%v", state.currentStep))
		}
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[superstep] unhandled corrdinator command: command=%#v", cmd))
		return
	}
}
func (state *coordinatorActor) computing(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.ComputeWorkerAck:
		if ok := state.ackRecorder.Ack(cmd.WorkerPid.GetId()); !ok {
			state.ActorUtil.LogError(context, fmt.Sprintf("compute ack from unknown worker: %v", cmd.WorkerPid))
			return
		}

		if cmd.AggregatedValues != nil {
			if err := aggregateValueMap(state.plugin.GetAggregators(), state.aggregatedCurrentStep, cmd.AggregatedValues); err != nil {
				state.ActorUtil.Fail(context, err)
				return
			}
		}
		if state.ackRecorder.HasCompleted() {
			state.ackRecorder.Clear()

			// check if there are active vertices
			stats, err := state.getStats(state.aggregatedCurrentStep)
			if err != nil {
				state.ActorUtil.Fail(context, err)
				return
			}

			// Why I check the number of messages sent is that the number of actives is often incorrect.
			// Vertex actor returns its active state with ComputeAck, but then it may receives a message until the next superstep is started
			if stats.ActiveVertices == 0 && stats.MessagesSent == 0 {
				// finish superstep
				state.behavior.Become(state.idle)
				state.ActorUtil.LogInfo(context, fmt.Sprintf("finish computing: step=%v", state.currentStep))

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
				state.ActorUtil.LogDebug(context, fmt.Sprintf("start computing: step=%v", state.currentStep))
			}

			// update aggregated values
			state.lastAggregatedValue.superstep = state.currentStep
			state.lastAggregatedValue.values = state.aggregatedCurrentStep
			state.aggregatedCurrentStep = make(map[string]*types.Any)
		}
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[computing] unhandled corrdinator command: command=%#v", cmd))
		return
	}
}

func (state *coordinatorActor) getStats(aggregated map[string]*types.Any) (*aggregator.VertexStats, error) {
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

func (state *coordinatorActor) findWorkerInfoByVertex(context actor.Context, vid plugin.VertexID) *command.ClusterInfo_WorkerInfo {
	p, err := state.plugin.Partition(vid, state.clusterInfo.NumOfPartitions())
	if err != nil {
		state.ActorUtil.LogError(context, fmt.Sprintf("failed to Partition(): %v", err))
		return nil
	}
	return state.clusterInfo.FindWoerkerInfoByPartition(p)
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
