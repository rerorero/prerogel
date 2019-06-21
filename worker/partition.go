package worker

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus"
)

type partitionActor struct {
	util.ActorUtil
	partitionID           uint64
	behavior              actor.Behavior
	plugin                plugin.Plugin
	vertices              map[plugin.VertexID]*actor.PID
	vertexProps           *actor.Props
	ackRecorder           *util.AckRecorder
	aggregatedCurrentStep map[string]*types.Any
}

// NewPartitionActor returns an actor instance
func NewPartitionActor(plg plugin.Plugin, vertexProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &util.AckRecorder{}
	ar.Clear()
	a := &partitionActor{
		plugin: plg,
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		vertexProps: vertexProps,
		vertices:    make(map[plugin.VertexID]*actor.PID),
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

		context.Respond(&command.InitPartitionAck{
			PartitionId: state.partitionID,
		})
		state.resetAckRecorder()
		state.behavior.Become(state.idle)
		state.ActorUtil.LogInfo(context, "initialization partition has completed")
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[waitInit] unhandled partition command: command=%#v", cmd))
		return
	}
}

func (state *partitionActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.LoadVertex:
		vid := plugin.VertexID(cmd.VertexId)
		if _, ok := state.vertices[vid]; ok {
			state.ActorUtil.LogError(context, fmt.Sprintf("vertex=%v has already created", cmd.VertexId))
		}
		pid, err := context.SpawnNamed(state.vertexProps, fmt.Sprintf("v%v", vid))
		if err != nil {
			state.ActorUtil.Fail(context, errors.Wrapf(err, "failed to spawn vertex %v", vid))
			return
		}
		state.vertices[vid] = pid
		context.Request(pid, cmd)
		return

	case *command.LoadVertexAck:
		context.Send(context.Parent(), cmd)
		return

	case *command.SuperStepBarrier:
		if len(state.vertices) == 0 {
			state.ActorUtil.LogInfo(context, fmt.Sprintf("[idle] no vertex is assigned"))
			return
		}
		state.resetAckRecorder()
		state.broadcastToVertices(context, cmd)
		state.behavior.Become(state.waitSuperStepBarrierAck)
		state.aggregatedCurrentStep = make(map[string]*types.Any)
		return
	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[idle] unhandled partition command: command=%#v", cmd))
		return
	}
}

func (state *partitionActor) waitSuperStepBarrierAck(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrierAck: // sent from parent
		if !state.ackRecorder.Ack(cmd.VertexId) {
			state.ActorUtil.LogWarn(context, fmt.Sprintf("SuperStepBarrierAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.HasCompleted() {
			context.Send(context.Parent(), &command.SuperStepBarrierPartitionAck{
				PartitionId: state.partitionID,
			})
			state.resetAckRecorder()
			state.behavior.Become(state.superstep)
			state.ActorUtil.LogInfo(context, "super step barrier has completed for partition")
		}
		return
	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[waitSpuerStepBarrierAck] unhandled partition command: command=%#v", cmd))
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
		// TODO: aggregate halted status
		if cmd.AggregatedValues != nil {
			if err := aggregateValueMap(state.plugin.GetAggregators(), state.aggregatedCurrentStep, cmd.AggregatedValues); err != nil {
				state.ActorUtil.Fail(context, err)
				return
			}
		}

		if !state.ackRecorder.Ack(cmd.VertexId) {
			state.ActorUtil.LogWarn(context, fmt.Sprintf("ComputeAck duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.HasCompleted() {
			context.Send(context.Parent(), &command.ComputePartitionAck{
				PartitionId:      state.partitionID,
				AggregatedValues: state.aggregatedCurrentStep,
			})
			state.resetAckRecorder()
			state.aggregatedCurrentStep = nil
			state.behavior.Become(state.idle)
			state.ActorUtil.LogInfo(context, "compute has completed for partition")
		}
		return

	case *command.SuperStepMessage:
		if _, ok := state.vertices[plugin.VertexID(cmd.SrcVertexId)]; ok {
			context.Forward(context.Parent())
		} else if pid, ok := state.vertices[plugin.VertexID(cmd.DestVertexId)]; ok {
			context.Forward(pid)
		} else {
			state.ActorUtil.LogError(context, fmt.Sprintf("[superstep] unknown destination message: msg=%#v", cmd))
		}
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[superstep] unhandled partition command: command=%#v", cmd))
		return
	}
}

func (state *partitionActor) broadcastToVertices(context actor.Context, msg interface{}) {
	state.LogDebug(context, fmt.Sprintf("broadcast %#v", msg))
	for _, pid := range state.vertices {
		context.Request(pid, msg)
	}
}

func (state *partitionActor) resetAckRecorder() {
	state.ackRecorder.Clear()
	for id := range state.vertices {
		state.ackRecorder.AddToWaitList(string(id))
	}
}
