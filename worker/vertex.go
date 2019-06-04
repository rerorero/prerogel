package worker

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus"
)

type vertexActor struct {
	plugin           Plugin
	vertex           Vertex
	logger           *logrus.Logger
	halted           bool
	receivedMessages []Message
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
	return c.vertexActor.receivedMessages
}

func (c *computeContextImpl) SendMessageTo(dest VertexID, m Message) error {
	msg, err := types.MarshalAny(m)
	if err != nil {
		c.vertexActor.logInfo(fmt.Sprintf("failed to marshal message: %v", err))
		return err
	}
	c.ctx.Send(c.ctx.Parent(), &command.SuperStepMessage{
		Uuid:         uuid.New().String(),
		SuperStep:    c.superStep,
		SrcVertexId:  string(c.vertexActor.vertex.GetID()),
		DestVertexId: string(dest),
		Message:      msg,
	})
	return nil
}

func (c *computeContextImpl) VoteToHalt() {
	c.vertexActor.logInfo(fmt.Sprintf("vertex %v has halted", c.vertexActor.vertex.GetID()))
	c.vertexActor.halted = true
}

// NewVertexActor returns an actor instance
func NewVertexActor(plugin Plugin, logger *logrus.Logger) actor.Actor {
	return &vertexActor{
		plugin: plugin,
		logger: logger,
	}
}

// Receive is message handler
func (state *vertexActor) Receive(context actor.Context) {
	switch cmd := context.Message().(type) {

	case *command.InitVertex:
		if state.vertex != nil {
			state.fail(errors.New("vertex has already initialized"))
			return
		}
		state.vertex = state.plugin.NewVertex(VertexID(cmd.VertexId))
		state.receivedMessages = nil
		state.logger = state.logger.WithField("vertexId", state.vertex.GetID()).Logger
		state.logInfo("vertex has initialized")
		context.Respond(&command.InitVertexAck{
			VertexId: string(state.vertex.GetID()),
		})

	case *command.LoadVertex:
		if state.vertex == nil {
			state.fail(fmt.Errorf("not initialized: cmd=%+v", cmd))
			return
		}
		if err := state.vertex.Load(); err != nil {
			state.fail(errors.New("failed to load vertex"))
			return
		}
		context.Respond(&command.LoadVertexAck{
			VertexId: string(state.vertex.GetID()),
		})

	case *command.Compute:
		if state.vertex == nil {
			state.fail(fmt.Errorf("not initialized: cmd=%+v", cmd))
			return
		}
		state.onComputed(context, cmd)

	case *command.SuperStepMessage:
		if state.vertex == nil {
			state.fail(fmt.Errorf("not initialized: cmd=%+v", cmd))
			return
		}
		if state.vertex.GetID() != VertexID(cmd.DestVertexId) {
			state.fail(fmt.Errorf("inconsistent vertex id: %+v", *cmd))
			return
		}
		// TODO: verify cmd.SuperStep
		state.receivedMessages = append(state.receivedMessages, cmd.Message)
		state.halted = false
		context.Respond(&command.SuperStepMessageAck{
			Uuid: cmd.Uuid,
		})

	case *actor.Started:
		state.logInfo("actor started")
	case *actor.Stopping:
		state.logInfo("actor stopping")
	case *actor.Stopped:
		state.logInfo("actor stopped")
	case *actor.Restarting:
		state.logInfo("actor restarting")
	default:
		state.fail(fmt.Errorf("unknown vertex command: command=%+v", cmd))
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
		state.logInfo("vertex is halted")
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
		state.fail(errors.Wrap(err, "failed to compute"))
		return
	}

	// clear messages
	state.receivedMessages = nil

	ctx.Respond(&command.ComputeAck{
		VertexId: string(id),
		Halted:   state.halted,
	})
	return
}

func (state *vertexActor) fail(err error) {
	state.logger.WithError(err).Error(err.Error())
	// let it crash
	panic(err)
}

func (state *vertexActor) logError(msg string) {
	state.logger.Error(msg)
}

func (state *vertexActor) logInfo(msg string) {
	state.logger.Info(msg)
}
