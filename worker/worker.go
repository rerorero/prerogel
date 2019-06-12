package worker

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/rerorero/prerogel/util"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus"
)

type superStepMsgBuf struct {
	buf    map[VertexID][]*command.SuperStepMessage
	plugin Plugin
}

type workerActor struct {
	util.ActorUtil
	behavior            actor.Behavior
	plugin              Plugin
	partitions          map[uint64]*actor.PID
	partitionProps      *actor.Props
	clusterInfo         *command.ClusterInfo
	ackRecorder         *ackRecorder
	combinedMessagesAck *ackRecorder
	ssMessageBuf        *superStepMsgBuf
}

// NewWorkerActor returns a new actor instance
func NewWorkerActor(plugin Plugin, partitionProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &ackRecorder{}
	ar.clear()
	mar := &ackRecorder{}
	mar.clear()
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
		info := workerInfoOf(cmd.ClusterInfo, context.Self())
		if info == nil {
			state.ActorUtil.Fail(fmt.Errorf("worker not assigned: me=%v", context.Self().GetId()))
			return
		}
		state.ActorUtil.AppendLoggerField("worker_id", context.Self().GetId())
		state.clusterInfo = cmd.ClusterInfo

		for _, partition := range info.Partitions {
			if _, ok := state.partitions[partition]; ok {
				state.ActorUtil.LogWarn(fmt.Sprintf("partition=%v has already created", partition))
				continue
			}

			pid := context.Spawn(state.partitionProps)
			state.partitions[partition] = pid
			context.Send(pid, &command.InitPartition{
				PartitionId: partition,
			})
		}
		state.resetAckRecorder()
		state.behavior.Become(state.waitPartitionInitAck)
		state.ActorUtil.LogDebug("become waitPartitionInitAck")
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled worker command: command=%+v", cmd))
		return
	}
}

func (state *workerActor) waitPartitionInitAck(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitPartitionAck:
		if !state.ackRecorder.ack(strconv.FormatUint(cmd.PartitionId, 10)) {
			state.ActorUtil.LogWarn(fmt.Sprintf("InitPartitionAck duplicated: id=%v", cmd.PartitionId))
		}
		if state.ackRecorder.hasCompleted() {
			context.Send(context.Parent(), &command.InitWorkerAck{
				WorkerPid: context.Self(),
			})
			state.resetAckRecorder()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogDebug("become superstep")
		}
		return
	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitPartitionInitÅck] unhandled worker command: command=%+v", cmd))
		return
	}
}

func (state *workerActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrier:
		state.ssMessageBuf.clear()
		state.broadcastToPartitions(context, cmd)
		state.resetAckRecorder()
		state.combinedMessagesAck.clear()
		state.behavior.Become(state.superstep)
		return
	default:
		state.ActorUtil.Fail(fmt.Errorf("[idle] unhandled worker command: command=%+v", cmd))
		return
	}
}

func (state *workerActor) waitSuperStepBarrierAck(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrierPartitionAck:
		if !state.ackRecorder.ack(strconv.FormatUint(cmd.PartitionId, 10)) {
			state.ActorUtil.LogWarn(fmt.Sprintf("SuperStepBarrierAck duplicated: id=%v", cmd.PartitionId))
		}
		if state.ackRecorder.hasCompleted() {
			context.Send(context.Parent(), &command.SuperStepBarrierWorkerAck{
				WorkerPid: context.Self(),
			})
			state.resetAckRecorder()
			state.behavior.Become(state.superstep)
			state.ActorUtil.LogInfo("super step barrier has completed for worker")
		}
		return
	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitSpuerStepBarrierAck] unhandled partition command: command=%+v", cmd))
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
		if !state.ackRecorder.ack(strconv.FormatUint(cmd.PartitionId, 10)) {
			state.ActorUtil.LogWarn(fmt.Sprintf("ComputeAck duplicated: id=%v", cmd.PartitionId))
		}
		if state.ackRecorder.hasCompleted() {
			if state.ssMessageBuf.numOfMessage() > 0 {
				if err := state.ssMessageBuf.combine(); err != nil {
					state.ActorUtil.LogError(fmt.Sprintf("failed to combine: %v", err))
				}
				// todo
				//for dest, msgs := range state.ssMessageBuf.buf {
				//	state.findWorkerInfoByPartition()
				//}
			} else {
				// TODO:
				//context.Send(context.Parent(), &command.ComputeWorkerAck{
				//	WorkerPid: context.Self(),
				//})
				//state.resetAckRecorder()
				//state.behavior.Become(state.idle)
				//state.ActorUtil.LogInfo("compute has completed for worker")
			}
		}
		// TODO: aggregate halted status
		return

	case *command.SuperStepMessage:
		srcWorker := state.findWorkerInfoByPartition(cmd.SrcPartitionId)
		if srcWorker == nil {
			state.ActorUtil.LogError(fmt.Sprintf("[superstep] message from unknown worker: command=%+v", cmd))
			return
		}

		if srcWorker.WorkerPid.GetId() == context.Self().GetId() {
			// when sent from my partition, saves to buffer then responds ack to vertex
			state.ssMessageBuf.add(cmd)
			if cmd.SrcVertexPid == nil {
				state.ActorUtil.Fail(fmt.Errorf("received message having invalid SrcertexPid: command=%+v", cmd))
				return
			}
			context.Send(cmd.SrcVertexPid, &command.SuperStepMessageAck{
				Uuid: cmd.Uuid,
			})
		} else {
			// when sent from other worker, route it to vertex
			destWorker := state.findWorkerInfoByVertex(VertexID(cmd.DestVertexId))
			if destWorker == nil {
				state.ActorUtil.LogError(fmt.Sprintf("[superstep] destination worker is not found: command=%+v", cmd))
				return
			}
			context.Send(pid, cmd)
		}
		return

	case *command.SuperStepMessageAck:
		//if state.combinedMessagesAck.ack(cmd.Uuid) {
		//	// original message was sent from my partition
		//	if state.combinedMessagesAck.hasCompleted() {
		//
		//	}
		//}
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[superstep] unhandled worker command: command=%+v", cmd))
		return
	}
}

func workerInfoOf(clusterInfo *command.ClusterInfo, pid *actor.PID) *command.WorkerInfo {
	for _, info := range clusterInfo.WorkerInfo {
		if info.WorkerPid.Id == pid.GetId() {
			return info
		}
	}
	return nil
}

func (state *workerActor) numOfPartitions() uint64 {
	var size uint64
	for _, i := range state.clusterInfo.WorkerInfo {
		size += uint64(len(i.Partitions))
	}
	return size
}

func (state *workerActor) broadcastToPartitions(context actor.Context, msg proto.Message) {
	state.LogDebug(fmt.Sprintf("broadcast %+v", msg))
	for _, pid := range state.partitions {
		context.Send(pid, msg)
	}
}

func (state *workerActor) resetAckRecorder() {
	state.ackRecorder.clear()
	for id := range state.partitions {
		state.ackRecorder.addToWaitList(strconv.FormatUint(id, 10))
	}
}

func (state *workerActor) findWorkerInfoByVertex(vid VertexID) *command.WorkerInfo {
	p, err := state.plugin.Partition(vid, state.numOfPartitions())
	if err != nil {
		state.ActorUtil.LogError(fmt.Sprintf("failed to Partition(): %v", err))
		return nil
	}
	return state.findWorkerInfoByPartition(p)
}

func (state *workerActor) findWorkerInfoByPartition(partitionID uint64) *command.WorkerInfo {
	for _, info := range state.clusterInfo.WorkerInfo {
		for _, p := range info.Partitions {
			if p == partitionID {
				return info
			}
		}
	}
	return nil
}

func (state *workerActor) computeAckAndBecomeIdle(context actor.Context) {
	context.Send(context.Parent(), &command.ComputeWorkerAck{
		WorkerPid: context.Self(),
	})
	state.resetAckRecorder()
	state.behavior.Become(state.idle)
	state.ActorUtil.LogInfo("compute has completed for worker")
}

// newSuperStepMsgBuf creates a new super step message buffer instance
func newSuperStepMsgBuf(plugin Plugin) *superStepMsgBuf {
	return &superStepMsgBuf{
		buf:    make(map[VertexID][]*command.SuperStepMessage),
		plugin: plugin,
	}
}

func (buf *superStepMsgBuf) clear() {
	buf.buf = make(map[VertexID][]*command.SuperStepMessage)
}

func (buf *superStepMsgBuf) numOfMessage() int {
	l := 0
	for _, s := range buf.buf {
		l += len(s)
	}
	return l
}

func (buf *superStepMsgBuf) add(m *command.SuperStepMessage) {
	buf.buf[VertexID(m.DestVertexId)] = append(buf.buf[VertexID(m.DestVertexId)], m)
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

		var msgs []Message
		for _, ss := range ssMsgs {
			m, err := buf.plugin.UnmarshalMessage(ss.Message)
			if err != nil {
				return errors.Wrapf(err, "failed to unmarshal message: %+v", ss)
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
				return errors.Wrapf(err, "failed to marshal combined message: %+v", c)
			}
			newMsgs = append(newMsgs, &command.SuperStepMessage{
				Uuid:           uuid.New().String(),
				SuperStep:      ssMsgs[0].SuperStep,
				SrcVertexId:    "",
				SrcPartitionId: ssMsgs[0].SrcPartitionId,
				SrcVertexPid:   nil,
				DestVertexId:   string(dest),
				Message:        pb,
			})
		}

		buf.buf[dest] = newMsgs
	}
	return nil
}
