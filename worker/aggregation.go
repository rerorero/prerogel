package worker

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

func aggregateValueMap(aggregators []Aggregator, base map[string]*any.Any, extra map[string]*any.Any) error {
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

		merged, err := agg.MarshalValue(agg.Aggregate(baseValue, extraValue))
		if err != nil {
			return errors.Wrapf(err, "failed to marshal aggregatable value: %#v, %#v", baseAny, extraAny)
		}

		base[name] = merged
	}

	return nil
}

func findAggregator(aggregators []Aggregator, name string) (Aggregator, error) {
	for _, a := range aggregators {
		if a.Name() == name {
			return a, nil
		}
	}
	return nil, fmt.Errorf("%s: no such aggregator", name)
}
