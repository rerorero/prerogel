package worker

import (
	"github.com/gogo/protobuf/types"
	"github.com/rerorero/prerogel/plugin"
)

type pluginProxy struct {
	underlying  plugin.Plugin
	aggregators []plugin.Aggregator
}

func newPluginProxy(p plugin.Plugin) *pluginProxy {
	return &pluginProxy{
		underlying:  p,
		aggregators: p.GetAggregators(),
	}
}

func (pp *pluginProxy) appendAggregators(agg []plugin.Aggregator) *pluginProxy {
	pp.aggregators = append(pp.aggregators, agg...)
	return pp
}

func (pp *pluginProxy) NewVertex(id plugin.VertexID) (plugin.Vertex, error) {
	return pp.underlying.NewVertex(id)
}

func (pp *pluginProxy) Partition(vertex plugin.VertexID, numOfPartitions uint64) (uint64, error) {
	return pp.underlying.Partition(vertex, numOfPartitions)
}

func (pp *pluginProxy) MarshalMessage(msg plugin.Message) (*types.Any, error) {
	return pp.underlying.MarshalMessage(msg)
}

func (pp *pluginProxy) UnmarshalMessage(pb *types.Any) (plugin.Message, error) {
	return pp.underlying.UnmarshalMessage(pb)
}

func (pp *pluginProxy) GetCombiner() func(destination plugin.VertexID, messages []plugin.Message) ([]plugin.Message, error) {
	return pp.underlying.GetCombiner()
}

func (pp *pluginProxy) GetAggregators() []plugin.Aggregator {
	return pp.aggregators
}
