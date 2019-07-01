package main

import (
	"fmt"
	"math"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/examples/sssp/loader"
	sssp "github.com/rerorero/prerogel/examples/sssp/proto"
	"github.com/rerorero/prerogel/plugin"
)

var _ = (plugin.Vertex)(&ssspVert{})
var _ = (plugin.Plugin)(&ssspPlugin{})

type ssspVert struct {
	id        string
	value     uint32
	outgoings map[string]uint32
	parent    string
}

func (v *ssspVert) Compute(ctx plugin.ComputeContext) error {
	defer ctx.VoteToHalt()

	if ctx.SuperStep() == 0 && v.value == 0 {
		// source vertex at super step 0
		// force it to send messages
	} else {
		min, err := getMinFromMessages(ctx.ReceivedMessages())
		if err != nil {
			return err
		}

		if min == nil || v.value < min.Value {
			return nil
		}

		v.value = min.Value
		v.parent = min.FromVertexId
	}

	for id, dist := range v.outgoings {
		if err := ctx.SendMessageTo(plugin.VertexID(id), &sssp.SSSPMessage{
			FromVertexId: v.id,
			Value:        v.value + dist,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (v *ssspVert) GetID() plugin.VertexID {
	return plugin.VertexID(v.id)
}

// ssspPlugin is single source shortest path plugin
type ssspPlugin struct {
	sourceID string
	graph    loader.Loader
}

func (p *ssspPlugin) NewVertex(id plugin.VertexID) (plugin.Vertex, error) {
	return nil, errors.New("NewVertex() is not implemented, you can load all vertices once using load command")
}

func (p *ssspPlugin) NewPartitionVertices(partitionID uint64, numOfPartitions uint64, register func(v plugin.Vertex)) error {
	values, err := p.graph.LoadPartition(partitionID, numOfPartitions)
	if err != nil {
		return err
	}

	for _, v := range values {
		initVal := uint32(math.MaxUint32)
		if v.ID == p.sourceID {
			initVal = 0
		}
		register(&ssspVert{
			id:        v.ID,
			value:     initVal,
			outgoings: v.Outgoings,
			parent:    "", // set none at the beginning
		})
	}

	return nil
}

// Partition provides partition information
func (p *ssspPlugin) Partition(vertex plugin.VertexID, numOfPartitions uint64) (uint64, error) {
	return plugin.HashPartition(vertex, numOfPartitions)
}

// MarshalMessage converts plugin.message to protobuf Any
func (p *ssspPlugin) MarshalMessage(msg plugin.Message) (*types.Any, error) {
	m, ok := msg.(*sssp.SSSPMessage)
	if !ok {
		return nil, fmt.Errorf("unexpected emssage type: %#v", msg)
	}

	return types.MarshalAny(m)
}

// UnmarshalMessage converts protobuf Any to plugin.message
func (p *ssspPlugin) UnmarshalMessage(pb *types.Any) (plugin.Message, error) {
	var m sssp.SSSPMessage
	if err := types.UnmarshalAny(pb, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func combiner(destination plugin.VertexID, messages []plugin.Message) ([]plugin.Message, error) {
	// choose the lowest value
	minMsg, err := getMinFromMessages(messages)
	if err != nil {
		return nil, err
	}
	return []plugin.Message{minMsg}, nil
}

func getMinFromMessages(messages []plugin.Message) (*sssp.SSSPMessage, error) {
	min := uint32(math.MaxUint32)
	var ret *sssp.SSSPMessage
	for _, m := range messages {
		msg, ok := m.(*sssp.SSSPMessage)
		if !ok {
			return nil, fmt.Errorf("unknown mesage type: %#v", m)
		}
		if msg.Value < min {
			ret = msg
		}
	}
	return ret, nil
}

// GetCombiner returns combiner function
func (p *ssspPlugin) GetCombiner() func(destination plugin.VertexID, messages []plugin.Message) ([]plugin.Message, error) {
	return combiner
}

// GetAggregators returns aggregators to be registered
func (p *ssspPlugin) GetAggregators() []plugin.Aggregator {
	return nil
}
