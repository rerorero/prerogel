package worker

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/util"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus"
)

type partitionActor struct {
	util.ActorUtil
	partitionID uint64
	behavior    actor.Behavior
	plugin      Plugin
	vertices    map[VertexID]*actor.PID
	vertexProps *actor.Props
	ackRecorder *ackRecorder
}

// NewPartitionActor returns an actor instance
func NewPartitionActor(plugin Plugin, vertexProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &ackRecorder{}
	ar.clear()
	a := &partitionActor{
		plugin: plugin,
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		vertexProps: vertexProps,
		vertices:    make(map[VertexID]*actor.PID),
		ackRecorder: ar,
	}
	a.behavior.Become(a.waitInit)
	return a
}

// Receive is message handler
func (state *partitionActor) Receive(context actor.Context) {
	if state.ActorUtil.IsSystemMessage(context.Message()) {
		// ignore
		return
	}
	state.behavior.Receive(context)
}

func (state *partitionActor) waitInit(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitPartition: // sent from parent
		state.partitionID = cmd.PartitionId
		state.ActorUtil.AppendLoggerField("partitionId", cmd.PartitionId)

		ids, err := state.plugin.ListVertexID(cmd.PartitionId)
		if err != nil {
			state.ActorUtil.Fail(errors.Wrap(err, "failed to ListVertexID()"))
			return
		}
		for _, id := range ids {
			if _, ok := state.vertices[id]; ok {
				state.ActorUtil.LogWarn(fmt.Sprintf("vertiex=%v has already created", id))
				continue
			}

			pid := context.Spawn(state.vertexProps)
			state.vertices[VertexID(id)] = pid
			context.Send(pid, &command.InitVertex{
				VertexId:    string(id),
				PartitionId: state.partitionID,
			})
		}
		state.resetAckRecorder()
		state.behavior.Become(state.waitInitVertexAck)
		state.ActorUtil.LogDebug("become waitInitVertexAck")
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled partition command: command=%+v", cmd))
		return
	}
}

func (state *partitionActor) waitInitVertexAck(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitVertexAck: // sent from parent
		if state.ackRecorder.ack(cmd.VertexId) {
			state.ActorUtil.LogWarn(fmt.Sprintf("InitVertexAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.hasCompleted() {
			context.Send(context.Parent(), &command.InitPartitionAck{
				PartitionId: state.partitionID,
			})
			state.resetAckRecorder()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogInfo("initialization partition has completed")
		}
	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInitVertexAck] unhandled partition command: command=%+v", cmd))
		return
	}
}

func (state *partitionActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrier:
		state.resetAckRecorder()
		state.broadcastToVertices(context, cmd)
		state.behavior.Become(state.waitSuperStepBarrierAck)
		return
	default:
		state.ActorUtil.Fail(fmt.Errorf("[idle] unhandled partition command: command=%+v", cmd))
		return
	}
}

func (state *partitionActor) waitSuperStepBarrierAck(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrierAck: // sent from parent
		if state.ackRecorder.ack(cmd.VertexId) {
			state.ActorUtil.LogWarn(fmt.Sprintf("SuperStepBarrierAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.hasCompleted() {
			context.Send(context.Parent(), &command.SuperStepBarrierPartitionAck{
				PartitionId: state.partitionID,
			})
			state.resetAckRecorder()
			state.behavior.Become(state.superstep)
			state.ActorUtil.LogInfo("super step barrier has completed for partition")
		}
		return
	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitSpuerStepBarrierAck] unhandled partition command: command=%+v", cmd))
		return
	}
}

func (state *partitionActor) superstep(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.Compute: // sent from parent
		state.resetAckRecorder()
		state.broadcastToVertices(context, cmd)
		return

	case *command.ComputeAck: // sent from vertices
		if state.ackRecorder.ack(cmd.VertexId) {
			state.ActorUtil.LogWarn(fmt.Sprintf("ComputeAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.hasCompleted() {
			context.Send(context.Parent(), &command.ComputePartitionAck{
				PartitionId: state.partitionID,
			})
			state.resetAckRecorder()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogInfo("compute has completed for partition")
		}
		// TODO: aggregate halted status
		return

	case *command.SuperStepMessage: // sent from vertices
		context.Forward(context.Parent())
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[superstep] unhandled partition command: command=%+v", cmd))
		return
	}
}

func (state *partitionActor) broadcastToVertices(context actor.Context, msg interface{}) {
	state.LogDebug(fmt.Sprintf("broadcast %+v", msg))
	for _, pid := range state.vertices {
		context.Send(pid, msg)
	}
}

func (state *partitionActor) resetAckRecorder() {
	state.ackRecorder.clear()
	for id := range state.vertices {
		state.ackRecorder.addToWaitList(string(id))
	}
}
