package worker

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
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
	ackRecorder *util.AckRecorder
}

// NewPartitionActor returns an actor instance
func NewPartitionActor(plugin Plugin, vertexProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &util.AckRecorder{}
	ar.Clear()
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

		context.Send(context.Parent(), &command.InitPartitionAck{
			PartitionId: state.partitionID,
		})
		state.resetAckRecorder()
		state.behavior.Become(state.idle)
		state.ActorUtil.LogInfo("initialization partition has completed")
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled partition command: command=%+v", cmd))
		return
	}
}

func (state *partitionActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.LoadVertex:
		vid := VertexID(cmd.VertexId)
		if _, ok := state.vertices[vid]; ok {
			state.ActorUtil.LogError(fmt.Sprintf("vertex=%v has already created", cmd.VertexId))
		}
		pid := context.Spawn(state.vertexProps)
		state.vertices[vid] = pid
		context.Send(pid, cmd)
		return

	case *command.LoadVertexAck:
		context.Send(context.Parent(), cmd)
		return

	case *command.SuperStepBarrier:
		if len(state.vertices) == 0 {
			state.ActorUtil.LogInfo(fmt.Sprintf("[idle] no vertex is assigned"))
			return
		}
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
		if !state.ackRecorder.Ack(cmd.VertexId) {
			state.ActorUtil.LogWarn(fmt.Sprintf("SuperStepBarrierAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.HasCompleted() {
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
		if !state.ackRecorder.Ack(cmd.VertexId) {
			state.ActorUtil.LogWarn(fmt.Sprintf("ComputeAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.HasCompleted() {
			context.Send(context.Parent(), &command.ComputePartitionAck{
				PartitionId: state.partitionID,
			})
			state.resetAckRecorder()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogInfo("compute has completed for partition")
		}
		// TODO: aggregate halted status
		return

	case *command.SuperStepMessage:
		if _, ok := state.vertices[VertexID(cmd.SrcVertexId)]; ok {
			context.Forward(context.Parent())
		} else if pid, ok := state.vertices[VertexID(cmd.DestVertexId)]; ok {
			context.Forward(pid)
		} else {
			state.ActorUtil.LogError(fmt.Sprintf("[superstep] unknown destination message: msg=%+v", cmd))
		}
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
	state.ackRecorder.Clear()
	for id := range state.vertices {
		state.ackRecorder.AddToWaitList(string(id))
	}
}
