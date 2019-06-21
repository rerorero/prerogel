package worker

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/plugin"
)

func aggregateValueMap(aggregators []plugin.Aggregator, base map[string]*any.Any, extra map[string]*any.Any) error {
	// TODO: there is still room for improvement
	for name, extraAny := range extra {
		baseAny, ok := base[name]
		if !ok {
			base[name] = extraAny
			continue
		}

		agg, err := findAggregator(aggregators, name)
		if err != nil {
			return err
		}

		baseValue, err := agg.UnmarshalValue(baseAny)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal aggregated value: %+v", baseAny)
		}

		extraValue, err := agg.UnmarshalValue(extraAny)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal aggregated value: %+v", extraAny)
		}

		val, err := agg.Aggregate(baseValue, extraValue)
		if err != nil {
			return errors.Wrapf(err, "failed to Aggregate()")
		}

		merged, err := agg.MarshalValue(val)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal aggregatable value: %#v, %#v", baseAny, extraAny)
		}

		base[name] = merged
	}

	return nil
}

func findAggregator(aggregators []plugin.Aggregator, name string) (plugin.Aggregator, error) {
	for _, a := range aggregators {
		if a.Name() == name {
			return a, nil
		}
	}
	return nil, fmt.Errorf("%s: no such aggregator", name)
}

func getAggregatedValue(aggregators []plugin.Aggregator, aggregated map[string]*any.Any, name string) (plugin.AggregatableValue, error) {
	ag, err := findAggregator(aggregators, name)
	if err != nil {
		return nil, err
	}

	anyVal, ok := aggregated[name]
	if !ok {
		return nil, fmt.Errorf("%s no such aggregated value", name)
	}

	return ag.UnmarshalValue(anyVal)
}
