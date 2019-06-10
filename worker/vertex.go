package worker

import (
	"fmt"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/util"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus"
)

type vertexActor struct {
	util.ActorUtil
	behavior         actor.Behavior
	plugin           Plugin
	vertex           Vertex
	partitionID      uint64
	halted           bool
	prevStepMessages []Message
	messageQueue     []Message
	ackRecorder      *ackRecorder
	computeRespondTo *actor.PID
}

type computeContextImpl struct {
	superStep   uint64
	ctx         actor.Context
	vertexActor *vertexActor
}

func (c *computeContextImpl) SuperStep() uint64 {
	return c.superStep
}

func (c *computeContextImpl) ReceivedMessages() []Message {
	return c.vertexActor.prevStepMessages
}

func (c *computeContextImpl) SendMessageTo(dest VertexID, m Message) error {
	msg, err := types.MarshalAny(m)
	if err != nil {
		c.vertexActor.ActorUtil.LogError(fmt.Sprintf("failed to marshal message: %v", err))
		return err
	}
	messageID := uuid.New().String()
	c.ctx.Send(c.ctx.Parent(), &command.SuperStepMessage{
		Uuid:           messageID,
		SuperStep:      c.superStep,
		SrcVertexId:    string(c.vertexActor.vertex.GetID()),
		SrcPartitionId: c.vertexActor.partitionID,
		SrcVertexPid:   c.ctx.Self(),
		DestVertexId:   string(dest),
		Message:        msg,
	})

	if c.vertexActor.ackRecorder.addToWaitList(messageID) {
		c.vertexActor.ActorUtil.LogWarn(fmt.Sprintf("duplicate superstep message: from=%v to=%v", c.vertexActor.vertex.GetID(), dest))
	}

	return nil
}

func (c *computeContextImpl) VoteToHalt() {
	c.vertexActor.ActorUtil.LogDebug("halted by user")
	c.vertexActor.halted = true
}

// NewVertexActor returns an actor instance
func NewVertexActor(plugin Plugin, logger *logrus.Logger) actor.Actor {
	ar := &ackRecorder{}
	ar.clear()
	a := &vertexActor{
		plugin: plugin,
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		ackRecorder: ar,
	}
	a.behavior.Become(a.waitInit)
	return a
}

// Receive is message handler
func (state *vertexActor) Receive(context actor.Context) {
	if state.ActorUtil.IsSystemMessage(context.Message()) {
		return
	}
	state.behavior.Receive(context)
}

func (state *vertexActor) waitInit(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitVertex:
		if state.vertex != nil {
			state.ActorUtil.Fail(errors.New("vertex has already initialized"))
			return
		}
		state.vertex = state.plugin.NewVertex(VertexID(cmd.VertexId))
		if err := state.vertex.Load(); err != nil {
			state.ActorUtil.Fail(errors.New("failed to load vertex"))
			return
		}
		state.partitionID = cmd.PartitionId
		state.ActorUtil.AppendLoggerField("vertexId", state.vertex.GetID())
		context.Respond(&command.InitVertexAck{
			VertexId: string(state.vertex.GetID()),
		})
		state.behavior.Become(state.superstep)
		state.ActorUtil.LogDebug("vertex has initialized")
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled vertex command: command=%+v(%v)", cmd, reflect.TypeOf(cmd)))
		return
	}
}

func (state *vertexActor) superstep(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrier:
		// move messages from queue to buffer
		state.prevStepMessages = state.messageQueue
		state.messageQueue = nil
		context.Respond(&command.SuperStepBarrierAck{
			VertexId: string(state.vertex.GetID()),
		})
		state.ActorUtil.LogDebug(fmt.Sprintf("received barrier message"))
		return

	case *command.Compute:
		state.computeRespondTo = context.Sender()
		state.onComputed(context, cmd)
		return

	case *command.SuperStepMessage:
		if state.vertex.GetID() != VertexID(cmd.DestVertexId) {
			state.ActorUtil.Fail(fmt.Errorf("inconsistent vertex id: %+v", *cmd))
			return
		}
		// TODO: verify cmd.SuperStep
		state.messageQueue = append(state.messageQueue, cmd.Message)
		state.halted = false
		context.Respond(&command.SuperStepMessageAck{
			Uuid: cmd.Uuid,
		})
		return

	case *command.SuperStepMessageAck:
		if state.ackRecorder.hasCompleted() {
			state.ActorUtil.LogWarn(fmt.Sprintf("unhaneled message id=%v, compute() has already completed", cmd.Uuid))
			return
		}
		if state.ackRecorder.ack(cmd.Uuid) {
			state.ActorUtil.LogWarn(fmt.Sprintf("duplicated or unhaneled message: uuid=%v", cmd.Uuid))
		}
		if state.ackRecorder.hasCompleted() {
			state.respondComputeAck(context)
		}
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[superstep] unhandled vertex command: command=%+v(%v)", cmd, reflect.TypeOf(cmd)))
		return
	}
}

func (state *vertexActor) onComputed(ctx actor.Context, cmd *command.Compute) {
	// force to compute() in super step 0
	// otherwise halt if there are no messages
	if cmd.SuperStep == 0 {
		state.halted = false
	} else if len(state.prevStepMessages) == 0 {
		state.ActorUtil.LogDebug("halted due to no message")
		state.halted = true
	}

	if state.halted {
		state.respondComputeAck(ctx)
		return
	}

	state.ackRecorder.clear()
	computeContext := &computeContextImpl{
		superStep:   cmd.SuperStep,
		ctx:         ctx,
		vertexActor: state,
	}
	if err := state.vertex.Compute(computeContext); err != nil {
		state.ActorUtil.Fail(errors.Wrap(err, "failed to compute"))
		return
	}

	if state.ackRecorder.hasCompleted() {
		state.respondComputeAck(ctx)
	}
	return
}

func (state *vertexActor) respondComputeAck(ctx actor.Context) {
	ctx.Send(state.computeRespondTo, &command.ComputeAck{
		VertexId: string(state.vertex.GetID()),
		Halted:   state.halted,
	})
	state.ActorUtil.LogDebug("compute() completed")
}
