package maximum

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/rerorero/prerogel/plugin"
)

// Plg is plugin that calculates maximum value of a graph
type Plg struct {
}

var _ = (plugin.Plugin)(&Plg{})

// NewVertex returns a new Vertex instance
func (p *Plg) NewVertex(id plugin.VertexID) plugin.Vertex {
	return nil
}

// Partition provides partition information
func (p *Plg) Partition(vertex plugin.VertexID, numOfPartitions uint64) (uint64, error) {
	return 0, nil
}

// MarshalMessage converts plugin.message to protobuf Any
func (p *Plg) MarshalMessage(msg plugin.Message) (*any.Any, error) {
	return nil, nil
}

// UnmarshalMessage converts protobuf Any to plugin.message
func (p *Plg) UnmarshalMessage(pb *any.Any) (plugin.Message, error) {
	return nil, nil
}

// GetCombiner returns combiner function
func (p *Plg) GetCombiner() func(destination plugin.VertexID, messages []plugin.Message) ([]plugin.Message, error) {
	return nil
}

// GetAggregators returns aggregators to be registered
func (p *Plg) GetAggregators() []plugin.Aggregator {
	return nil
}
