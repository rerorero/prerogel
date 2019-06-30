package plugin

import (
	"github.com/gogo/protobuf/types"
)

// VertexValue indicates value which a vertex holds.
type VertexValue interface{}

// EdgeValue indicates value which an edge holds.
type EdgeValue interface{}

// Message is a message sent from the vertex to another vertex during super-step.
type Message interface{}

// VertexID is id of vertex
type VertexID string

// AggregatableValue is value to be aggregated by aggregator
type AggregatableValue interface{}

// ComputeContext provides information for vertices to process Compute()
type ComputeContext interface {
	SuperStep() uint64
	ReceivedMessages() []Message
	SendMessageTo(dest VertexID, m Message) error
	VoteToHalt()
	GetAggregated(aggregatorName string) (AggregatableValue, bool, error)
	PutAggregatable(aggregatorName string, v AggregatableValue) error
}

// Vertex is abstract of a vertex. thread safe.
type Vertex interface {
	Compute(computeContext ComputeContext) error
	GetID() VertexID
}

// Aggregator is Pregel aggregator implemented by user
type Aggregator interface {
	Name() string
	Aggregate(v1 AggregatableValue, v2 AggregatableValue) (AggregatableValue, error)
	MarshalValue(v AggregatableValue) (*types.Any, error)
	UnmarshalValue(pb *types.Any) (AggregatableValue, error)
	ToString(v AggregatableValue) string
}

// Plugin is a plugin that provides graph computation.
type Plugin interface {
	NewVertex(id VertexID) (Vertex, error)
	Partition(vertex VertexID, numOfPartitions uint64) (uint64, error)
	NewPartitionVertices(partitionID uint64, numOfPartitions uint64, register func(v Vertex)) error
	MarshalMessage(msg Message) (*types.Any, error)
	UnmarshalMessage(pb *types.Any) (Message, error)
	GetCombiner() func(destination VertexID, messages []Message) ([]Message, error)
	GetAggregators() []Aggregator
}
