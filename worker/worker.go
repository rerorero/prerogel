package worker

import (
	"fmt"
	"strconv"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/rerorero/prerogel/util"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus"
)

type workerActor struct {
	util.ActorUtil
	behavior       actor.Behavior
	plugin         Plugin
	partitions     map[uint64]*actor.PID
	partitionProps *actor.Props
	clusterInfo    *command.ClusterInfo
	ackRecorder    *ackRecorder
	ssMessageBuf   map[VertexID][]*command.SuperStepMessage
}

// NewWorkerActor returns a new actor instance
func NewWorkerActor(plugin Plugin, partitionProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &ackRecorder{}
	ar.clear()
	a := &workerActor{
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		plugin:         plugin,
		partitions:     make(map[uint64]*actor.PID),
		partitionProps: partitionProps,
		ackRecorder:    ar,
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
		info, ok := workerInfoOf(cmd.ClusterInfo, context.Self())
		if !ok {
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
		if state.ackRecorder.ack(strconv.FormatUint(cmd.PartitionId, 10)) {
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
		state.clearSuperStepMessageBuff()
		state.broadcastToPartitions(context, cmd)
		state.resetAckRecorder()
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
		if state.ackRecorder.ack(strconv.FormatUint(cmd.PartitionId, 10)) {
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
		if state.ackRecorder.ack(strconv.FormatUint(cmd.PartitionId, 10)) {
			state.ActorUtil.LogWarn(fmt.Sprintf("ComputeAck duplicated: id=%v", cmd.PartitionId))
		}
		if state.ackRecorder.hasCompleted() {
			context.Send(context.Parent(), &command.ComputeWorkerAck{
				WorkerPid: context.Self(),
			})
			state.resetAckRecorder()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogInfo("compute has completed for worker")
		}
		// TODO: aggregate halted status
		return

	case *command.SuperStepMessage:
		if _, ok := state.partitions[cmd.SrcPartitionId]; ok {
			// when sent from my partition, saves to buffer and respond ack to vertex
			state.ssMessageBuf[VertexID(cmd.DestVertexId)] = append(state.ssMessageBuf[VertexID(cmd.DestVertexId)], cmd)
			context.Send(cmd.SrcVertexPid, &command.SuperStepMessageAck{
				Uuid: cmd.Uuid,
			})
		} else {
			// when sent from other worker, route it to vertex

		}
		return
	default:
		state.ActorUtil.Fail(fmt.Errorf("[superstep] unhandled worker command: command=%+v", cmd))
		return
	}
}

func workerInfoOf(clusterInfo *command.ClusterInfo, pid *actor.PID) (*command.WorkerInfo, bool) {
	for _, info := range clusterInfo.WorkerInfo {
		if info.WorkerPid.Id == pid.GetId() {
			return info, true
		}
	}
	return nil, false
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

func (state *workerActor) clearSuperStepMessageBuff() {
	state.ssMessageBuf = make(map[VertexID][]*command.SuperStepMessage)
}
