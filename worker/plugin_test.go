package worker

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/rerorero/prerogel/plugin"
)

// MockedPlugin is mocked Plugin struct
type MockedPlugin struct {
	NewVertexMock        func(id plugin.VertexID) plugin.Vertex
	PartitionMock        func(id plugin.VertexID, numOfPartitions uint64) (uint64, error)
	MarshalMessageMock   func(msg plugin.Message) (*any.Any, error)
	UnmarshalMessageMock func(a *any.Any) (plugin.Message, error)
	GetCombinerMock      func() func(plugin.VertexID, []plugin.Message) ([]plugin.Message, error)
	GetAggregatorsMock   func() []plugin.Aggregator
}

func (m *MockedPlugin) NewVertex(id plugin.VertexID) plugin.Vertex {
	return m.NewVertexMock(id)
}

func (m *MockedPlugin) Partition(id plugin.VertexID, numOfPartitions uint64) (uint64, error) {
	return m.PartitionMock(id, numOfPartitions)
}

func (m *MockedPlugin) MarshalMessage(msg plugin.Message) (*any.Any, error) {
	return m.MarshalMessageMock(msg)
}

func (m *MockedPlugin) UnmarshalMessage(a *any.Any) (plugin.Message, error) {
	return m.UnmarshalMessageMock(a)
}

func (m *MockedPlugin) GetCombiner() func(destination plugin.VertexID, messages []plugin.Message) ([]plugin.Message, error) {
	return m.GetCombinerMock()
}

func (m *MockedPlugin) GetAggregators() []plugin.Aggregator {
	return m.GetAggregatorsMock()
}

// MockedVertex is mocked Vertex struct
type MockedVertex struct {
	LoadMock        func() error
	ComputeMock     func(computeContext plugin.ComputeContext) error
	GetIDMock       func() plugin.VertexID
	GetOutEdgesMock func() []plugin.Edge
	GetValueMock    func() (plugin.VertexValue, error)
	SetValueMock    func(v plugin.VertexValue) error
}

func (m *MockedVertex) Load() error {
	return m.LoadMock()
}
func (m *MockedVertex) Compute(computeContext plugin.ComputeContext) error {
	return m.ComputeMock(computeContext)
}
func (m *MockedVertex) GetID() plugin.VertexID {
	return m.GetIDMock()
}
func (m *MockedVertex) GetOutEdges() []plugin.Edge {
	return m.GetOutEdgesMock()
}
func (m *MockedVertex) GetValue() (plugin.VertexValue, error) {
	return m.GetValueMock()
}
func (m *MockedVertex) SetValue(v plugin.VertexValue) error {
	return m.SetValueMock(v)
}

// MockedAggregator
type MockedAggregator struct {
	NameMock           func() string
	AggregateMock      func(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) plugin.AggregatableValue
	MarshalValueMock   func(v plugin.AggregatableValue) (*any.Any, error)
	UnmarshalValueMock func(pb *any.Any) (plugin.AggregatableValue, error)
}

func (m *MockedAggregator) Name() string {
	return m.NameMock()
}

func (m *MockedAggregator) Aggregate(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) plugin.AggregatableValue {
	return m.AggregateMock(v1, v2)
}

func (m *MockedAggregator) MarshalValue(v plugin.AggregatableValue) (*any.Any, error) {
	return m.MarshalValueMock(v)

}

func (m *MockedAggregator) UnmarshalValue(pb *any.Any) (plugin.AggregatableValue, error) {
	return m.UnmarshalValueMock(pb)
}
