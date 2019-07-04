package worker

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
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

	switch cmd := context.Message().(type) {
	case *command.ClusterInfo:
		state.broadcastToVertices(context, cmd)
		return

	case *command.GetVertexValue:
		if v, ok := state.vertices[plugin.VertexID(cmd.VertexId)]; ok {
			context.Forward(v)
		} else {
			state.LogWarn(context, fmt.Sprintf("%v no such vertex id in partition %v", cmd.VertexId, state.partitionID))
			context.Respond(&command.GetVertexValueAck{VertexId: cmd.VertexId})
		}
		return

	default:
		state.behavior.Receive(context)
		return
	}
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
			err := fmt.Sprintf("vertex has already created: id=%s", cmd.VertexId)
			state.ActorUtil.LogError(context, err)
			context.Respond(&command.LoadVertexAck{VertexId: string(cmd.VertexId), Error: err})
			return
		}
		pid, err := context.SpawnNamed(state.vertexProps, fmt.Sprintf("v%v", vid))
		if err != nil {
			err := fmt.Sprintf("failed to spawn actor: id=%s", cmd.VertexId)
			state.ActorUtil.LogError(context, err)
			context.Respond(&command.LoadVertexAck{VertexId: string(cmd.VertexId), Error: err})
			return
		}
		state.vertices[vid] = pid
		context.Forward(pid)
		return

	case *command.LoadPartitionVertices:
		state.resetAckRecorder()
		var loadErr string
		if err := state.plugin.NewPartitionVertices(state.partitionID, cmd.NumOfPartitions, func(v plugin.Vertex) {
			// TODO: concurrency unsafe
			vid := v.GetID()
			pid, err := context.SpawnNamed(state.vertexProps, fmt.Sprintf("v%v", vid))
			if err != nil {
				loadErr = fmt.Sprintf("failed to spawn actor: id=%s", vid)
				state.ActorUtil.LogError(context, loadErr)
				return
			}
			state.vertices[vid] = pid
			context.Request(pid, &loadVertexLocal{vertex: v})
			state.ackRecorder.AddToWaitList(string(vid))
		}); err != nil {
			state.ackRecorder.Clear()
			state.ActorUtil.LogError(context, err.Error())
			context.Respond(&command.LoadPartitionVerticesAck{PartitionId: state.partitionID, Error: err.Error()})
			return
		}
		if loadErr != "" {
			state.ackRecorder.Clear()
			context.Respond(&command.LoadPartitionVerticesAck{PartitionId: state.partitionID, Error: loadErr})
			return
		}
		state.ActorUtil.LogDebug(context, fmt.Sprintf("start waiting for loading partition vertices"))
		state.behavior.Become(state.waitLoadPartitionVertices)
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

	case *command.SuperStepMessage:
		state.handleMessage(context, cmd)
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[idle] unhandled partition command: command=%#v", cmd))
		return
	}
}

func (state *partitionActor) waitLoadPartitionVertices(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.LoadVertexAck:
		if !state.ackRecorder.Ack(cmd.VertexId) {
			state.ActorUtil.LogWarn(context, fmt.Sprintf("LoadVertexAck(partition) duplicated: id=%v", cmd.VertexId))
		}
		if state.ackRecorder.HasCompleted() {
			context.Send(context.Parent(), &command.LoadPartitionVerticesAck{
				PartitionId: state.partitionID,
			})
			state.resetAckRecorder()
			state.behavior.Become(state.idle)
			state.ActorUtil.LogDebug(context, "loading partition finished")
		}
		return
	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[waitLoadPartitionVertices] unhandled partition command: command=%#v", cmd))
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
			state.ActorUtil.LogDebug(context, "partition: super step barrier end")
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
			state.ActorUtil.LogInfo(context, "partition: compute has completed")
		}
		return

	case *command.SuperStepMessage:
		state.handleMessage(context, cmd)
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[superstep] unhandled partition command: command=%#v", cmd))
		return
	}
}
func (state *partitionActor) handleMessage(context actor.Context, cmd *command.SuperStepMessage) {
	if _, ok := state.vertices[plugin.VertexID(cmd.SrcVertexId)]; ok {
		context.Forward(context.Parent())
	} else if pid, ok := state.vertices[plugin.VertexID(cmd.DestVertexId)]; ok {
		context.Forward(pid)
	} else {
		state.ActorUtil.LogError(context, fmt.Sprintf("[superstep] unknown destination message: msg=%#v", cmd))
	}
}

func (state *partitionActor) broadcastToVertices(context actor.Context, msg interface{}) {
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
