package aggregator

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/rerorero/prerogel/plugin"
)

// VertexStatsName is aggregator name of VertexStatsAggregator
const VertexStatsName = "prerogel/vertex-stats"

// VertexStatsAggregatorInstance is singleton
var VertexStatsAggregatorInstance = &VertexStatsAggregator{}

// VertexStatsAggregator is aggregator for VertexStats
type VertexStatsAggregator struct{}

// Name is aggregator name
func (s *VertexStatsAggregator) Name() string {
	return VertexStatsName
}

// Aggregate is reduction func
func (s *VertexStatsAggregator) Aggregate(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) plugin.AggregatableValue {
	pb1, ok := v1.(*VertexStats)
	if !ok {
		return fmt.Errorf("unknown aggregatable value: %#v", v1)
	}
	pb2, ok := v2.(*VertexStats)
	if !ok {
		return fmt.Errorf("unknown aggregatable value: %#v", v2)
	}
	return &VertexStats{
		ActiveVertices: pb1.ActiveVertices + pb2.ActiveVertices,
		TotalVertices:  pb1.TotalVertices + pb2.TotalVertices,
		MessagesSent:   pb1.MessagesSent + pb2.MessagesSent,
	}
}

// MarshalValue converts AggregatableValue into any.Any
func (s *VertexStatsAggregator) MarshalValue(v plugin.AggregatableValue) (*any.Any, error) {
	pb, ok := v.(*VertexStats)
	if !ok {
		return nil, fmt.Errorf("unknown aggregatable value: %#v", v)
	}
	return ptypes.MarshalAny(pb)
}

// UnmarshalValue converts any.Any to AggregatableValue
func (s *VertexStatsAggregator) UnmarshalValue(pb *any.Any) (plugin.AggregatableValue, error) {
	var stats VertexStats
	if err := ptypes.UnmarshalAny(pb, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
