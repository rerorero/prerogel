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
	partitionID    uint64
	nrOfPartitions uint64
	behavior       actor.Behavior
	plugin         Plugin
	vertices       map[VertexID]*actor.PID
	vertexProps    *actor.Props
	ackRecorder    *ackRecorder
	superStep      uint64
}

type ackRecorder struct {
	m map[VertexID]struct{}
}

// NewPartitionActor returns an actor instance
func NewPartitionActor(plugin Plugin, logger *logrus.Logger) actor.Actor {
	ar := &ackRecorder{}
	ar.clear()
	a := &partitionActor{
		plugin: plugin,
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		vertexProps: actor.PropsFromProducer(func() actor.Actor {
			return NewVertexActor(plugin, logger)
		}),
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
	case *command.InitPartition:
		state.partitionID = cmd.PartitionId
		state.nrOfPartitions = cmd.NrOfPartitions
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
				VertexId: string(id),
			})
		}
		state.ackRecorder.clear()
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
	case *command.InitVertexAck:
		if state.ackRecorder.ack(VertexID(cmd.VertexId)) {
			state.ActorUtil.LogWarn(fmt.Sprintf("InitVertexAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.hasCompleted(state.vertices) {
			context.Send(context.Parent(), &command.InitPartitionAck{
				PartitionId: state.partitionID,
			})
			state.ackRecorder.clear()
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
	case *command.LoadPartition:
		state.broadcastToVertices(context, &command.LoadVertex{})
		state.behavior.Become(state.waitLoadVertexAck)
		return
	case *command.Compute:
		state.superStep = cmd.SuperStep
		state.broadcastToVertices(context, cmd)
		state.behavior.Become(state.computing)
		return
	default:
		state.ActorUtil.Fail(fmt.Errorf("[idle] unhandled partition command: command=%+v", cmd))
		return
	}
}

func (state *partitionActor) waitLoadVertexAck(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.LoadVertexAck:
		if state.ackRecorder.ack(VertexID(cmd.VertexId)) {
			state.ActorUtil.LogWarn(fmt.Sprintf("LoadVetexAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.hasCompleted(state.vertices) {
			context.Send(context.Parent(), &command.LoadPartitionAck{
				PartitionId: state.partitionID,
			})
			state.ackRecorder.clear()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogInfo("load partition has completed")
		}
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitLoadVertexAck] unhandled partition command: command=%+v", cmd))
		return
	}
}
func (state *partitionActor) computing(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.ComputeAck:
		if state.ackRecorder.ack(VertexID(cmd.VertexId)) {
			state.ActorUtil.LogWarn(fmt.Sprintf("ComputeAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.hasCompleted(state.vertices) {
			// TODO: return ack
		}
		// TODO: aggregate halted status
		return

	case *command.SuperStepMessage:
		context.Forward(context.Parent())
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[computing] unhandled partition command: command=%+v", cmd))
		return
	}
}

func (state *partitionActor) broadcastToVertices(context actor.Context, msg interface{}) {
	state.LogDebug(fmt.Sprintf("broadcast %+v", msg))
	for _, pid := range state.vertices {
		context.Send(pid, msg)
	}
}
