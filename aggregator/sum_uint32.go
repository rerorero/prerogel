package aggregator

import (
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/types"
	"github.com/rerorero/prerogel/plugin"
)

// SumUint32Aggregator is aggregator for SumUint32Aggregator
type SumUint32Aggregator struct {
	aggName string
}

// NewSumUint32Aggregator returns a new SumUint32Aggregator instance
func NewSumUint32Aggregator(name string) *SumUint32Aggregator {
	return &SumUint32Aggregator{
		aggName: name,
	}
}

// Name returns aggregator name
func (s *SumUint32Aggregator) Name() string {
	return s.aggName
}

// Aggregate is reduction func
func (s *SumUint32Aggregator) Aggregate(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) (plugin.AggregatableValue, error) {
	n1, ok := v1.(uint32)
	if !ok {
		return nil, fmt.Errorf("expected aggregatable value type is uint32: %#v", v1)
	}
	n2, ok := v2.(uint32)
	if !ok {
		return nil, fmt.Errorf("expected aggregatable value type is uint32: %#v", v2)
	}
	return n1 + n2, nil
}

// MarshalValue converts AggregatableValue into uin32
func (s *SumUint32Aggregator) MarshalValue(v plugin.AggregatableValue) (*types.Any, error) {
	return plugin.ConvertUint32ToAny(v)
}

// UnmarshalValue converts uint32 to AggregatableValue
func (s *SumUint32Aggregator) UnmarshalValue(pb *types.Any) (plugin.AggregatableValue, error) {
	return plugin.ConvertAnyToUint32(pb)
}

// ToString converts aggregatabale value to string
func (s *SumUint32Aggregator) ToString(v plugin.AggregatableValue) string {
	n, ok := v.(uint32)
	if !ok {
		return fmt.Sprintf("<unknown: not uint32 %v>", v)
	}
	return strconv.FormatUint(uint64(n), 10)
}
