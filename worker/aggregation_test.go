package worker

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/go-cmp/cmp"
	"github.com/rerorero/prerogel/plugin"
)

func Test_aggregateValueMap(t *testing.T) {
	type args struct {
		aggregators []plugin.Aggregator
		base        map[string]*any.Any
		extra       map[string]*any.Any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantBase map[string]*any.Any
	}{
		{
			name: "aggregate",
			args: args{
				aggregators: aggregators,
				base: map[string]*any.Any{
					"concat": {Value: []byte("AA")},
				},
				extra: map[string]*any.Any{
					"concat": {Value: []byte("BB")},
					"sum":    {Value: []byte{uint8(3)}},
				},
			},
			wantErr: false,
			wantBase: map[string]*any.Any{
				"concat": {Value: []byte("AABB")},
				"sum":    {Value: []byte{uint8(3)}},
			},
		},
		{
			name: "nil extra",
			args: args{
				aggregators: aggregators,
				base: map[string]*any.Any{
					"concat": {Value: []byte("AA")},
				},
				extra: nil,
			},
			wantErr: false,
			wantBase: map[string]*any.Any{
				"concat": {Value: []byte("AA")},
			},
		},
		{
			name: "unknown",
			args: args{
				aggregators: aggregators,
				base: map[string]*any.Any{
					"foo": {Value: []byte("AA")},
				},
				extra: map[string]*any.Any{
					"foo": {Value: []byte{uint8(3)}},
				},
			},
			wantErr: true,
			wantBase: map[string]*any.Any{
				"foo": {Value: []byte("AA")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := aggregateValueMap(tt.args.aggregators, tt.args.base, tt.args.extra); (err != nil) != tt.wantErr {
				t.Fatalf("aggregateValueMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		if diff := cmp.Diff(tt.args.base, tt.wantBase); diff != "" {
			t.Errorf("different base: %s", diff)
		}
	}
}
