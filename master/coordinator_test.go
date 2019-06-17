package master

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_assignPartition(t *testing.T) {
	type args struct {
		workers        int
		nrOfPartitions uint64
	}
	tests := []struct {
		name    string
		args    args
		want    [][]uint64
		wantErr bool
	}{
		{
			name: "worker > 1, partition > 1, surplus > 0",
			args: args{
				workers:        3,
				nrOfPartitions: 11,
			},
			want: [][]uint64{
				{0, 1, 2, 3},
				{4, 5, 6, 7},
				{8, 9, 10},
			},
			wantErr: false,
		},
		{
			name: "worker > 1, partition > 1, surplus = 0",
			args: args{
				workers:        3,
				nrOfPartitions: 9,
			},
			want: [][]uint64{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			wantErr: false,
		},
		{
			name: "worker = 1, partition > 1",
			args: args{
				workers:        1,
				nrOfPartitions: 3,
			},
			want: [][]uint64{
				{0, 1, 2},
			},
			wantErr: false,
		},
		{
			name: "worker > 1, partition = 1",
			args: args{
				workers:        3,
				nrOfPartitions: 1,
			},
			want: [][]uint64{
				{0},
				{},
				{},
			},
			wantErr: false,
		},
		{
			name: "worker = 0",
			args: args{
				workers:        0,
				nrOfPartitions: 10,
			},
			wantErr: true,
		},
		{
			name: "partitions = 0",
			args: args{
				workers:        2,
				nrOfPartitions: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assignPartition(tt.args.workers, tt.args.nrOfPartitions)
			if (err != nil) != tt.wantErr {
				t.Errorf("assignPartition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("assignPartition() = %s", diff)
			}
		})
	}
}
