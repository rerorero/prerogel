package worker

import (
	"github.com/golang/protobuf/ptypes/any"
)

// MockedPlugin is mocked Plugin struct
type MockedPlugin struct {
	NewVertexMock        func(id VertexID) Vertex
	PartitionMock        func(id VertexID, numOfPartitions uint64) (uint64, error)
	MarshalMessageMock   func(msg Message) (*any.Any, error)
	UnmarshalMessageMock func(a *any.Any) (Message, error)
	GetCombinerMock      func() func(VertexID, []Message) ([]Message, error)
	GetAggregatorsMock   func() []Aggregator
}

func (m *MockedPlugin) NewVertex(id VertexID) Vertex {
	return m.NewVertexMock(id)
}

func (m *MockedPlugin) Partition(id VertexID, numOfPartitions uint64) (uint64, error) {
	return m.PartitionMock(id, numOfPartitions)
}

func (m *MockedPlugin) MarshalMessage(msg Message) (*any.Any, error) {
	return m.MarshalMessageMock(msg)
}

func (m *MockedPlugin) UnmarshalMessage(a *any.Any) (Message, error) {
	return m.UnmarshalMessageMock(a)
}

func (m *MockedPlugin) GetCombiner() func(destination VertexID, messages []Message) ([]Message, error) {
	return m.GetCombinerMock()
}

func (m *MockedPlugin) GetAggregators() []Aggregator {
	return m.GetAggregatorsMock()
}

// MockedVertex is mocked Vertex struct
type MockedVertex struct {
	LoadMock        func() error
	ComputeMock     func(computeContext ComputeContext) error
	GetIDMock       func() VertexID
	GetOutEdgesMock func() []Edge
	GetValueMock    func() (VertexValue, error)
	SetValueMock    func(v VertexValue) error
}

func (m *MockedVertex) Load() error {
	return m.LoadMock()
}
func (m *MockedVertex) Compute(computeContext ComputeContext) error {
	return m.ComputeMock(computeContext)
}
func (m *MockedVertex) GetID() VertexID {
	return m.GetIDMock()
}
func (m *MockedVertex) GetOutEdges() []Edge {
	return m.GetOutEdgesMock()
}
func (m *MockedVertex) GetValue() (VertexValue, error) {
	return m.GetValueMock()
}
func (m *MockedVertex) SetValue(v VertexValue) error {
	return m.SetValueMock(v)
}

// MockedAggregator
type MockedAggregator struct {
	NameMock           func() string
	AggregateMock      func(v1 AggregatableValue, v2 AggregatableValue) AggregatableValue
	MarshalValueMock   func(v AggregatableValue) (*any.Any, error)
	UnmarshalValueMock func(pb *any.Any) (AggregatableValue, error)
}

func (m *MockedAggregator) Name() string {
	return m.NameMock()
}

func (m *MockedAggregator) Aggregate(v1 AggregatableValue, v2 AggregatableValue) AggregatableValue {
	return m.AggregateMock(v1, v2)
}

func (m *MockedAggregator) MarshalValue(v AggregatableValue) (*any.Any, error) {
	return m.MarshalValueMock(v)

}

func (m *MockedAggregator) UnmarshalValue(pb *any.Any) (AggregatableValue, error) {
	return m.UnmarshalValueMock(pb)
}
