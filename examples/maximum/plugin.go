package main

import (
	"fmt"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/examples/maximum/loader"
	"github.com/rerorero/prerogel/plugin"
)

const aggregatorName = "prerogel-maximum"

var _ = (plugin.Vertex)(&vert{})
var _ = (plugin.Plugin)(&maxPlugin{})
var aggregators = []plugin.Aggregator{&maxAggregator{}}

// vert is vertex of graph
type vert struct {
	id            string
	value         uint32
	outgoingEdges []string
}

func (v *vert) Compute(ctx plugin.ComputeContext) error {
	// send self to all outgoing edges
	for _, edge := range v.outgoingEdges {
		if err := ctx.SendMessageTo(plugin.VertexID(edge), v.value); err != nil {
			return err
		}
	}

	// no messages are sent in superstep 0
	if ctx.SuperStep() == 0 {
		return nil
	}

	messages := ctx.ReceivedMessages()

	if len(messages) == 0 {
		ctx.VoteToHalt()
		return nil
	}

	max, err := getMaxFromMessages(messages)
	if err != nil {
		return err
	}

	if max < v.value {
		ctx.VoteToHalt()
		return nil
	}

	return nil
}

func (v *vert) GetID() plugin.VertexID {
	return plugin.VertexID(v.id)
}

// maxAggregator is aggregator of max
type maxAggregator struct{}

func (agg *maxAggregator) Name() string {
	return aggregatorName
}

func (agg *maxAggregator) Aggregate(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) (plugin.AggregatableValue, error) {
	v1n, ok1 := v1.(uint32)
	v2n, ok2 := v2.(uint32)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("unknown value: v1=%#v, v2=%#v", v1, v2)
	}
	if v1n > v2n {
		return v1n, nil
	}
	return v2n, nil
}

func (agg *maxAggregator) MarshalValue(v plugin.AggregatableValue) (*types.Any, error) {
	return plugin.ConvertUint32ToAny(v)
}

func (agg *maxAggregator) UnmarshalValue(pb *types.Any) (plugin.AggregatableValue, error) {
	return plugin.ConvertAnyToUint32(pb)
}

// maxPlugin is maximum value plugin
type maxPlugin struct {
	graph loader.Loader
}

// NewVertex returns a new Vertex instance
func (p *maxPlugin) NewVertex(id plugin.VertexID) (plugin.Vertex, error) {
	value, outgoings, err := p.graph.Load(string(id))
	if err != nil {
		return nil, nil
	}
	return &vert{
		id:            string(id),
		value:         value,
		outgoingEdges: outgoings,
	}, nil
}

// Partition provides partition information
func (p *maxPlugin) Partition(vertex plugin.VertexID, numOfPartitions uint64) (uint64, error) {
	return plugin.HashPartition(vertex, numOfPartitions)
}

// MarshalMessage converts plugin.message to protobuf Any
func (p *maxPlugin) MarshalMessage(msg plugin.Message) (*types.Any, error) {
	return plugin.ConvertUint32ToAny(msg)
}

// UnmarshalMessage converts protobuf Any to plugin.message
func (p *maxPlugin) UnmarshalMessage(pb *types.Any) (plugin.Message, error) {
	return plugin.ConvertAnyToUint32(pb)
}

func combiner(destination plugin.VertexID, messages []plugin.Message) ([]plugin.Message, error) {
	// choose the highest number
	max, err := getMaxFromMessages(messages)
	if err != nil {
		return nil, err
	}
	return []plugin.Message{max}, nil
}

func getMaxFromMessages(messages []plugin.Message) (uint32, error) {
	if len(messages) == 0 {
		return 0, errors.New("expects non empty slice")
	}
	var max uint32
	for _, m := range messages {
		mv, ok := m.(uint32)
		if !ok {
			return 0, fmt.Errorf("unknown mesage type: %#v", m)
		}
		if mv > max {
			max = mv
		}
	}
	return max, nil
}

// GetCombiner returns combiner function
func (p *maxPlugin) GetCombiner() func(destination plugin.VertexID, messages []plugin.Message) ([]plugin.Message, error) {
	return combiner
}

// GetAggregators returns aggregators to be registered
func (p *maxPlugin) GetAggregators() []plugin.Aggregator {
	return aggregators
}
