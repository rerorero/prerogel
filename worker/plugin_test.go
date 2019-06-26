package worker

import (
	"github.com/gogo/protobuf/types"
	"github.com/rerorero/prerogel/plugin"
)

// MockedPlugin is mocked Plugin struct
type MockedPlugin struct {
	NewVertexMock        func(id plugin.VertexID) (plugin.Vertex, error)
	PartitionMock        func(id plugin.VertexID, numOfPartitions uint64) (uint64, error)
	MarshalMessageMock   func(msg plugin.Message) (*types.Any, error)
	UnmarshalMessageMock func(a *types.Any) (plugin.Message, error)
	GetCombinerMock      func() func(plugin.VertexID, []plugin.Message) ([]plugin.Message, error)
	GetAggregatorsMock   func() []plugin.Aggregator
}

func (m *MockedPlugin) NewVertex(id plugin.VertexID) (plugin.Vertex, error) {
	return m.NewVertexMock(id)
}

func (m *MockedPlugin) Partition(id plugin.VertexID, numOfPartitions uint64) (uint64, error) {
	return m.PartitionMock(id, numOfPartitions)
}

func (m *MockedPlugin) MarshalMessage(msg plugin.Message) (*types.Any, error) {
	return m.MarshalMessageMock(msg)
}

func (m *MockedPlugin) UnmarshalMessage(a *types.Any) (plugin.Message, error) {
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
	ComputeMock func(computeContext plugin.ComputeContext) error
	GetIDMock   func() plugin.VertexID
}

func (m *MockedVertex) Compute(computeContext plugin.ComputeContext) error {
	return m.ComputeMock(computeContext)
}
func (m *MockedVertex) GetID() plugin.VertexID {
	return m.GetIDMock()
}

// MockedAggregator
type MockedAggregator struct {
	NameMock           func() string
	AggregateMock      func(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) (plugin.AggregatableValue, error)
	MarshalValueMock   func(v plugin.AggregatableValue) (*types.Any, error)
	UnmarshalValueMock func(pb *types.Any) (plugin.AggregatableValue, error)
	ToStringMock       func(v plugin.AggregatableValue) string
}

func (m *MockedAggregator) Name() string {
	return m.NameMock()
}

func (m *MockedAggregator) Aggregate(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) (plugin.AggregatableValue, error) {
	return m.AggregateMock(v1, v2)
}

func (m *MockedAggregator) MarshalValue(v plugin.AggregatableValue) (*types.Any, error) {
	return m.MarshalValueMock(v)
}

func (m *MockedAggregator) UnmarshalValue(pb *types.Any) (plugin.AggregatableValue, error) {
	return m.UnmarshalValueMock(pb)
}

func (m *MockedAggregator) ToString(v plugin.AggregatableValue) string {
	return m.ToStringMock(v)
}
