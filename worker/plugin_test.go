package worker

import (
	"github.com/golang/protobuf/ptypes/any"
)

// MockedPlugin is mocked Plugin struct
type MockedPlugin struct {
	NewVertexMock        func(id VertexID) Vertex
	ListVertexIDMock     func(partitionId uint64) ([]VertexID, error)
	MarshalMessageMock   func(msg Message) (*any.Any, error)
	UnmarshalMessageMock func(a *any.Any) (Message, error)
	GetCombinerMock      func() func(VertexID, []Message) ([]Message, error)
}

func (m *MockedPlugin) NewVertex(id VertexID) Vertex {
	return m.NewVertexMock(id)
}

func (m *MockedPlugin) ListVertexID(partitionID uint64) ([]VertexID, error) {
	return m.ListVertexIDMock(partitionID)
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
