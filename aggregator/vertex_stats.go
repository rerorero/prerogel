package aggregator

import (
	"fmt"

	"github.com/gogo/protobuf/types"
	"github.com/rerorero/prerogel/plugin"
)

// VertexStatsAggregator is aggregator for VertexStats
type VertexStatsAggregator struct {
	AggName string
}

// Name is aggregator name
func (s *VertexStatsAggregator) Name() string {
	return s.AggName
}

// Aggregate is reduction func
func (s *VertexStatsAggregator) Aggregate(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) (plugin.AggregatableValue, error) {
	pb1, ok := v1.(*VertexStats)
	if !ok {
		return nil, fmt.Errorf("unknown aggregatable value: %#v", v1)
	}
	pb2, ok := v2.(*VertexStats)
	if !ok {
		return nil, fmt.Errorf("unknown aggregatable value: %#v", v2)
	}
	return &VertexStats{
		ActiveVertices: pb1.ActiveVertices + pb2.ActiveVertices,
		TotalVertices:  pb1.TotalVertices + pb2.TotalVertices,
		MessagesSent:   pb1.MessagesSent + pb2.MessagesSent,
	}, nil
}

// MarshalValue converts AggregatableValue into types.Any
func (s *VertexStatsAggregator) MarshalValue(v plugin.AggregatableValue) (*types.Any, error) {
	pb, ok := v.(*VertexStats)
	if !ok {
		return nil, fmt.Errorf("unknown aggregatable value: %#v", v)
	}
	return types.MarshalAny(pb)
}

// UnmarshalValue converts types.Any to AggregatableValue
func (s *VertexStatsAggregator) UnmarshalValue(pb *types.Any) (plugin.AggregatableValue, error) {
	var stats VertexStats
	if err := types.UnmarshalAny(pb, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
