package worker

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus"
)

type superStepMsgBuf struct {
	buf    map[plugin.VertexID][]*command.SuperStepMessage
	plugin plugin.Plugin
}

type workerActor struct {
	util.ActorUtil
	coordinatorPID        *actor.PID
	behavior              actor.Behavior
	plugin                plugin.Plugin
	partitions            map[uint64]*actor.PID
	partitionProps        *actor.Props
	clusterInfo           *command.ClusterInfo
	ackRecorder           *util.AckRecorder
	combinedMessagesAck   *util.AckRecorder
	ssMessageBuf          *superStepMsgBuf
	aggregatedCurrentStep map[string]*types.Any
}

// NewWorkerActor returns a new actor instance
func NewWorkerActor(plugin plugin.Plugin, partitionProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &util.AckRecorder{}
	ar.Clear()
	mar := &util.AckRecorder{}
	mar.Clear()
	a := &workerActor{
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		plugin:              plugin,
		partitions:          make(map[uint64]*actor.PID),
		partitionProps:      partitionProps,
		ackRecorder:         ar,
		combinedMessagesAck: mar,
		ssMessageBuf:        newSuperStepMsgBuf(plugin),
	}
	a.behavior.Become(a.waitInit)
	return a
}

// Receive is message handler
func (state *workerActor) Receive(context actor.Context) {
	if state.ActorUtil.IsSystemMessage(context.Message()) {
		// ignore
		return
	}
	state.behavior.Receive(context)
}

func (state *workerActor) waitInit(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitWorker:
		state.coordinatorPID = cmd.Coordinator
		for _, partition := range cmd.Partitions {
			if _, ok := state.partitions[partition]; ok {
				state.ActorUtil.LogWarn(context, fmt.Sprintf("partition=%v has already created", partition))
				continue
			}

			pid, error := context.SpawnNamed(state.partitionProps, fmt.Sprintf("p%v", partition))
			if error != nil {
				state.ActorUtil.Fail(context, errors.Wrapf(error, "spawn partition failed %v", partition))
				return
			}
			state.partitions[partition] = pid
			context.Request(pid, &command.InitPartition{
				PartitionId: partition,
			})
		}
		state.resetAckRecorder()
		state.behavior.Become(state.waitPartitionInitAck)
		state.ActorUtil.LogDebug(context, "become waitPartitionInitAck")
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[waitInit] unhandled worker command: command=%#v", cmd))
		return
	}
}

func (state *workerActor) waitPartitionInitAck(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitPartitionAck:
		state.ActorUtil.LogDebug(context, fmt.Sprintf("partitionInitAck from id=%v", cmd.PartitionId))
		if !state.ackRecorder.Ack(strconv.FormatUint(cmd.PartitionId, 10)) {
			state.ActorUtil.LogWarn(context, fmt.Sprintf("InitPartitionAck duplicated: id=%v", cmd.PartitionId))
		}
		if state.ackRecorder.HasCompleted() {
			context.Send(state.coordinatorPID, &command.InitWorkerAck{
				WorkerPid: context.Self(),
			})
			state.resetAckRecorder()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogDebug(context, "become idle")
		}
		return
	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[waitPartitionInitÅck] unhandled worker command: command=%#v", cmd))
		return
	}
}

func (state *workerActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.LoadVertex:
		destPartition, err := state.plugin.Partition(plugin.VertexID(cmd.VertexId), state.clusterInfo.NumOfPartitions())
		if err != nil {
			err := fmt.Sprintf("failed to find partition: vertex id=%s err=%v", cmd.VertexId, err)
			state.ActorUtil.LogError(context, err)
			context.Respond(&command.LoadVertexAck{VertexId: string(cmd.VertexId), Error: err})
			return
		}

		destPid, ok := state.partitions[destPartition]
		if !ok {
			err := fmt.Sprintf("routing vertex issue: vertex id=%s", cmd.VertexId)
			state.ActorUtil.LogError(context, err)
			context.Respond(&command.LoadVertexAck{VertexId: string(cmd.VertexId), Error: err})
			return
		}

		context.Forward(destPid)
		return

	case *command.SuperStepBarrier:
		state.ActorUtil.LogDebug(context, "super step barrier")
		if err := state.checkClusterInfo(cmd.ClusterInfo, context.Self()); err != nil {
			state.ActorUtil.Fail(context, err)
			return
		}
		state.clusterInfo = cmd.ClusterInfo
		state.clusterInfo = cmd.ClusterInfo
		state.ssMessageBuf.clear()
		state.broadcastToPartitions(context, cmd)
		state.resetAckRecorder()
		state.aggregatedCurrentStep = make(map[string]*types.Any)
		state.combinedMessagesAck.Clear()
		state.behavior.Become(state.waitSuperStepBarrierAck)
		state.ActorUtil.LogDebug(context, "become waitSuperStepBarrierAck")
		return

	case *command.SuperStepMessage:
		// need to handle because other workers might still run super-step
		state.handleSuperStepMessage(context, cmd)
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[idle] unhandled worker command: command=%#v", cmd))
		return
	}
}

func (state *workerActor) waitSuperStepBarrierAck(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrierPartitionAck:
		state.ActorUtil.LogDebug(context, fmt.Sprintf("super step barrier partition ack: id=%v", cmd.PartitionId))
		if !state.ackRecorder.Ack(strconv.FormatUint(cmd.PartitionId, 10)) {
			state.ActorUtil.LogWarn(context, fmt.Sprintf("SuperStepBarrierAck duplicated: id=%v", cmd.PartitionId))
		}
		if state.ackRecorder.HasCompleted() {
			context.Send(state.coordinatorPID, &command.SuperStepBarrierWorkerAck{
				WorkerPid: context.Self(),
			})
			state.resetAckRecorder()
			state.behavior.Become(state.superstep)
			state.ActorUtil.LogInfo(context, "super step barrier has completed for worker")
		}
		return
	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[waitSpuerStepBarrierAck] unhandled partition command: command=%#v", cmd))
		return
	}
}

func (state *workerActor) superstep(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.Compute: // sent from parent
		state.resetAckRecorder()
		state.broadcastToPartitions(context, cmd)
		return

	case *command.ComputePartitionAck:
		// TODO: aggregate halted status
		if cmd.AggregatedValues != nil {
			if err := aggregateValueMap(state.plugin.GetAggregators(), state.aggregatedCurrentStep, cmd.AggregatedValues); err != nil {
				state.ActorUtil.Fail(context, err)
				return
			}
		}
		if !state.ackRecorder.Ack(strconv.FormatUint(cmd.PartitionId, 10)) {
			state.ActorUtil.LogWarn(context, fmt.Sprintf("ComputeAck duplicated: id=%v", cmd.PartitionId))
		}
		if state.ackRecorder.HasCompleted() {
			if state.ssMessageBuf.numOfMessage() > 0 {
				if err := state.ssMessageBuf.combine(); err != nil {
					state.ActorUtil.LogError(context, fmt.Sprintf("failed to combine: %v", err))
				}
				// TODO: it can reduce messages by aggregating by each destination worker
				for dest, msgs := range state.ssMessageBuf.buf {
					destWorker := state.findWorkerInfoByVertex(context, dest)
					if destWorker == nil || destWorker.WorkerPid.GetId() == context.Self().GetId() {
						state.ActorUtil.Fail(context, fmt.Errorf("failed to find worker: %v", dest))
						return
					}
					for _, m := range msgs {
						context.Request(destWorker.WorkerPid, m)
					}
				}
				// wait for SuperStepMessageAck from other workers
			} else {
				state.computeAckAndBecomeIdle(context)
			}
		}
		return

	case *command.SuperStepMessage:
		state.handleSuperStepMessage(context, cmd)
		return

	case *command.SuperStepMessageAck:
		state.ssMessageBuf.remove(cmd)
		if state.ssMessageBuf.numOfMessage() == 0 {
			state.computeAckAndBecomeIdle(context)
		}

		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[superstep] unhandled worker command: command=%#v", cmd))
		return
	}
}

func (state *workerActor) handleSuperStepMessage(context actor.Context, cmd *command.SuperStepMessage) {
	srcWorker := state.findWorkerInfoByVertex(context, plugin.VertexID(cmd.SrcVertexId))
	if srcWorker == nil {
		state.ActorUtil.LogError(context, fmt.Sprintf("[superstep] message from unknown worker: command=%#v", cmd))
		return
	}

	if srcWorker.WorkerPid.GetId() == context.Self().GetId() {
		destPartition, err := state.plugin.Partition(plugin.VertexID(cmd.DestVertexId), state.clusterInfo.NumOfPartitions())
		if err != nil {
			state.ActorUtil.Fail(context, fmt.Errorf("failed to find partition for message: %#v", cmd))
			return
		}

		destPid, ok := state.partitions[destPartition]
		if ok {
			// when sent from local partition to local partition, forward it
			context.Forward(destPid)

		} else {
			// when sent from local partition to other worker's partition, saves to buffer then responds Ack to vertex
			state.ssMessageBuf.add(cmd)
			context.Respond(&command.SuperStepMessageAck{
				Uuid: cmd.Uuid,
			})
		}

	} else {
		// when sent from other worker, route it to vertex
		p, err := state.plugin.Partition(plugin.VertexID(cmd.DestVertexId), state.clusterInfo.NumOfPartitions())
		if err != nil {
			state.ActorUtil.Fail(context, errors.Wrap(err, "failed to Partition()"))
			return
		}
		pid, ok := state.partitions[p]
		if !ok {
			state.ActorUtil.Fail(context, fmt.Errorf("[superstep] destination partition(%v) is not found: command=%#v", p, cmd))
			return
		}
		context.Forward(pid)
	}
	return
}

func (state *workerActor) checkClusterInfo(clusterInfo *command.ClusterInfo, self *actor.PID) error {
	var cmdPartitions []uint64
	for _, i := range clusterInfo.WorkerInfo {
		if i.WorkerPid.Id == self.GetId() {
			cmdPartitions = i.Partitions
			break
		}
	}
	if cmdPartitions == nil {
		return fmt.Errorf("worker not found in command.WorkerInfo: me=%v", self)
	}
	var currentPartitions []uint64
	for p := range state.partitions {
		currentPartitions = append(currentPartitions, p)
	}
	sort.Slice(cmdPartitions, func(i, j int) bool { return cmdPartitions[i] < cmdPartitions[j] })
	sort.Slice(currentPartitions, func(i, j int) bool { return currentPartitions[i] < currentPartitions[j] })

	if !reflect.DeepEqual(cmdPartitions, currentPartitions) {
		// TODO: reconcile partitions
		return fmt.Errorf("assigned partitions has changed: %v -> %v", currentPartitions, cmdPartitions)
	}
	return nil
}

func (state *workerActor) broadcastToPartitions(context actor.Context, msg proto.Message) {
	state.LogDebug(context, fmt.Sprintf("broadcast %#v", msg))
	for _, pid := range state.partitions {
		context.Request(pid, msg)
	}
}

func (state *workerActor) resetAckRecorder() {
	state.ackRecorder.Clear()
	for id := range state.partitions {
		state.ackRecorder.AddToWaitList(strconv.FormatUint(id, 10))
	}
}

func (state *workerActor) findWorkerInfoByVertex(context actor.Context, vid plugin.VertexID) *command.ClusterInfo_WorkerInfo {
	p, err := state.plugin.Partition(vid, state.clusterInfo.NumOfPartitions())
	if err != nil {
		state.ActorUtil.LogError(context, fmt.Sprintf("failed to Partition(): %v", err))
		return nil
	}
	return state.clusterInfo.FindWoerkerInfoByPartition(p)
}

func (state *workerActor) computeAckAndBecomeIdle(context actor.Context) {
	context.Send(state.coordinatorPID, &command.ComputeWorkerAck{
		WorkerPid:        context.Self(),
		AggregatedValues: state.aggregatedCurrentStep,
	})
	state.aggregatedCurrentStep = nil
	state.resetAckRecorder()
	state.behavior.Become(state.idle)
	state.ActorUtil.LogInfo(context, "compute has completed for worker")
}

// newSuperStepMsgBuf creates a new super step message buffer instance
func newSuperStepMsgBuf(plg plugin.Plugin) *superStepMsgBuf {
	return &superStepMsgBuf{
		buf:    make(map[plugin.VertexID][]*command.SuperStepMessage),
		plugin: plg,
	}
}

func (buf *superStepMsgBuf) clear() {
	buf.buf = make(map[plugin.VertexID][]*command.SuperStepMessage)
}

func (buf *superStepMsgBuf) numOfMessage() int {
	l := 0
	for _, s := range buf.buf {
		l += len(s)
	}
	return l
}

func (buf *superStepMsgBuf) add(m *command.SuperStepMessage) {
	buf.buf[plugin.VertexID(m.DestVertexId)] = append(buf.buf[plugin.VertexID(m.DestVertexId)], m)
}

func (buf *superStepMsgBuf) remove(ack *command.SuperStepMessageAck) {
	for vid, msgs := range buf.buf {
		for i, m := range msgs {
			if m.Uuid == ack.Uuid {
				removed := append(msgs[:i], msgs[i+1:]...)
				if len(removed) == 0 {
					delete(buf.buf, vid)
				} else {
					buf.buf[vid] = removed
				}
				return
			}
		}
	}
}

func (buf *superStepMsgBuf) combine() error {
	combiner := buf.plugin.GetCombiner()
	if combiner == nil {
		return nil
	}

	for dest, ssMsgs := range buf.buf {
		if len(ssMsgs) <= 1 {
			continue
		}

		var msgs []plugin.Message
		for _, ss := range ssMsgs {
			m, err := buf.plugin.UnmarshalMessage(ss.Message)
			if err != nil {
				return errors.Wrapf(err, "failed to unmarshal message: %#v", ss)
			}
			msgs = append(msgs, m)
		}

		combined, err := combiner(dest, msgs)
		if err != nil {
			return errors.Wrapf(err, "failed to combine message: dest=%v", dest)
		}

		var newMsgs []*command.SuperStepMessage
		for _, c := range combined {
			pb, err := buf.plugin.MarshalMessage(c)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal combined message: %#v", c)
			}
			newMsgs = append(newMsgs, &command.SuperStepMessage{
				Uuid:         uuid.New().String(),
				SuperStep:    ssMsgs[0].SuperStep,
				SrcVertexId:  "",
				DestVertexId: string(dest),
				Message:      pb,
			})
		}

		buf.buf[dest] = newMsgs
	}
	return nil
}
