package pregel

import (
	"fmt"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pkg/errors"
	"github.com/rerorero/graph/pregel/actor/command"
)

type VertexValue interface{}
type EdgeValue interface{}
type Message interface{}
type VertexID string

type Edge struct {
	Value  EdgeValue
	Target VertexID
}

type ComputeContext interface {
	SuperStep() uint64
	ReceivedMessages() []Message
	SendMessageTo(dest VertexID, m Message)
	VoteToHalt()
}

// Vertex is abstract of vertex. used by only one channel (thread safe)
type Vertex interface {
	Load() error
	Compute(computeContext ComputeContext) error
	GetID() VertexID
	GetOutEdges() []Edge
	GetValue() (VertexValue, error)
	SetValue(v VertexValue) error
}

type Plugin interface {
	NewVertex(id VertexID) Vertex
}

type vertexActor struct {
	plugin           Plugin
	vertex           Vertex
	logger           *log.Logger
	halted           bool
	receivedMessages []Message
}

func NewVertexActor(plugin Plugin, logger *log.Logger) actor.Actor {
	return &vertexActor{
		logger: logger,
	}
}

// Receive is handler of Actor message
func (state *vertexActor) Receive(context actor.Context) {
	msg := context.Message()

	// initialize
	if initCommand, ok := msg.(*command.InitVertex); ok {
		if state.vertex != nil {
			state.fail(fmt.Errorf("vertex has already initialized with %v: received = %v", state.vertex.GetID(), initCommand.VertexId))
			return
		}
		state.vertex = state.plugin.NewVertex(VertexID(initCommand.VertexId))
		state.halted = false
		state.receivedMessages = nil
		return
	}

	// other commands
	if state.vertex == nil {
		state.fail(fmt.Errorf("not initialized: cmd=%+v", msg))
		return
	}
	id := state.vertex.GetID()

	switch cmd := msg.(type) {

	case *command.LoadVertex:
		if err := state.vertex.Load(); err != nil {
			state.fail(fmt.Errorf("failed to load vertex id=%v", id))
			return
		}

	case *command.Compute:
		if state.halted {
			state.logger.Printf("vertex id=%v is halted\n", id)
			return
		}
		ctx, err := state.getComputeContext(cmd)
		if err != nil {
			state.fail(err)
			return
		}
		if err := state.vertex.Compute(ctx); err != nil {
			state.fail(errors.Wrapf(err, "failed to compute: %v", id))
			return
		}
		// clear messages
		state.receivedMessages = nil

	case *command.CustomVertexMessage:
		state.receivedMessages = append(state.receivedMessages, cmd.Message)

	default:
		state.fail(fmt.Errorf("unknown vertex command: id=%v, command=%+v", id, cmd))
	}
}

func (state *vertexActor) getComputeContext(cmd *command.Compute) (ComputeContext, error) {
	// TODO
	return nil, nil
}

func (state *vertexActor) fail(err error) {
	state.logger.Println(err)
}
