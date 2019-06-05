package worker

import "github.com/gogo/protobuf/proto"

// VertexValue indicates value which a vertex holds.
type VertexValue interface{}

// EdgeValue indicates value which an edge holds.
type EdgeValue interface{}

// Message is a message sent from the vertex to another vertex during super-step.
type Message proto.Message

// VertexID is id of vertex
type VertexID string

// Edge indicates an edge of graph
type Edge struct {
	Value  EdgeValue
	Target VertexID
}

// ComputeContext provides information for vertices to process Compute()
type ComputeContext interface {
	SuperStep() uint64
	ReceivedMessages() []Message
	SendMessageTo(dest VertexID, m Message) error
	VoteToHalt()
}

// Vertex is abstract of a vertex. thread safe.
type Vertex interface {
	Load() error
	Compute(computeContext ComputeContext) error
	GetID() VertexID
	GetOutEdges() []Edge
	GetValue() (VertexValue, error)
	SetValue(v VertexValue) error
}

// Plugin is a plugin that provides graph computation.
type Plugin interface {
	NewVertex(id VertexID) Vertex
	ListVertexID(partitionID uint64) ([]VertexID, error)
}
