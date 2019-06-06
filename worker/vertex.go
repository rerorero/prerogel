package worker

import (
	"fmt"

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
	halted           bool
	receivedMessages []Message
	ackRecorder      *ackRecorder
}

type computeContextImpl struct {
	superStep        uint64
	ctx              actor.Context
	vertexActor      *vertexActor
	sentMessagesDest map[VertexID]struct{}
}

func (c *computeContextImpl) SuperStep() uint64 {
	return c.superStep
}

func (c *computeContextImpl) ReceivedMessages() []Message {
	return c.vertexActor.receivedMessages
}

func (c *computeContextImpl) SendMessageTo(dest VertexID, m Message) error {
	msg, err := types.MarshalAny(m)
	if err != nil {
		c.vertexActor.ActorUtil.LogError(fmt.Sprintf("failed to marshal message: %v", err))
		return err
	}
	c.ctx.Send(c.ctx.Parent(), &command.SuperStepMessage{
		Uuid:         uuid.New().String(),
		SuperStep:    c.superStep,
		SrcVertexId:  string(c.vertexActor.vertex.GetID()),
		DestVertexId: string(dest),
		Message:      msg,
	})

	if _, ok := c.sentMessagesDest[dest]; ok {
		c.vertexActor.ActorUtil.LogWarn(fmt.Sprintf("duplicate superstep message: from=%v to=%v", c.vertexActor.vertex.GetID(), dest))
	}
	c.sentMessagesDest[dest] = struct{}{}
	return nil
}

func (c *computeContextImpl) VoteToHalt() {
	c.vertexActor.ActorUtil.LogDebug(fmt.Sprintf("vertex %v has halted", c.vertexActor.vertex.GetID()))
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
}

func (state *vertexActor) waitInit(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitVertex:
		if state.vertex != nil {
			state.ActorUtil.Fail(errors.New("vertex has already initialized"))
			return
		}
		state.vertex = state.plugin.NewVertex(VertexID(cmd.VertexId))
		state.receivedMessages = nil
		state.ActorUtil.AppendLoggerField("vertexId", state.vertex.GetID())
		context.Respond(&command.InitVertexAck{
			VertexId: string(state.vertex.GetID()),
		})
		state.behavior.Become(state.idle)
		state.ActorUtil.LogDebug("vertex has initialized")
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled vertex command: command=%+v", cmd))
		return
	}
}

func (state *vertexActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.LoadVertex:
		if err := state.vertex.Load(); err != nil {
			state.ActorUtil.Fail(errors.New("failed to load vertex"))
			return
		}
		context.Respond(&command.LoadVertexAck{
			VertexId: string(state.vertex.GetID()),
		})
		return

	case *command.Compute:
		state.onComputed(context, cmd)
		return

	case *command.SuperStepMessage:
		if state.vertex.GetID() != VertexID(cmd.DestVertexId) {
			state.ActorUtil.Fail(fmt.Errorf("inconsistent vertex id: %+v", *cmd))
			return
		}
		// TODO: verify cmd.SuperStep
		state.receivedMessages = append(state.receivedMessages, cmd.Message)
		state.halted = false
		context.Respond(&command.SuperStepMessageAck{
			Uuid: cmd.Uuid,
		})
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[idle] unhandled vertex command: command=%+v", cmd))
		return
	}
}

func (state *vertexActor) onComputed(ctx actor.Context, cmd *command.Compute) {
	id := state.vertex.GetID()

	// always do compute() in super step 0
	// otherwise halt when not receiving messages
	if cmd.SuperStep == 0 {
		state.halted = false
	} else if len(state.receivedMessages) == 0 {
		state.halted = true
	}

	if state.halted {
		state.ActorUtil.LogDebug("vertex is halted")
		ctx.Respond(&command.ComputeAck{
			VertexId: string(id),
			Halted:   state.halted,
		})
		return
	}
	computeContext := &computeContextImpl{
		superStep:   cmd.SuperStep,
		ctx:         ctx,
		vertexActor: state,
	}
	if err := state.vertex.Compute(computeContext); err != nil {
		state.ActorUtil.Fail(errors.Wrap(err, "failed to compute"))
		return
	}

	// clear messages
	// TODO: sepalate message buffer (current step and previous step)
	state.receivedMessages = nil

	state.ackRecorder.clear()
	state.ackRecorder.hasCompleted(len(computeContext.sentMessagesDest))
	ctx.Respond(&command.ComputeAck{
		VertexId: string(id),
		Halted:   state.halted,
	})
	return
}
