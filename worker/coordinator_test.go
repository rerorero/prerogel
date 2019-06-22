package worker

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"github.com/google/go-cmp/cmp"
	"github.com/rerorero/prerogel/aggregator"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus/hooks/test"
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

func TestNewCoordinatorActor(t *testing.T) {
	var mux sync.Mutex
	var initCount int32
	var loadCount int32
	var barrierCount int32
	var receivedPartitions []uint64
	var stepCount int32
	waitCh := make(chan string, 1)
	logger, _ := test.NewNullLogger()
	plugin := &MockedPlugin{
		PartitionMock: func(id plugin.VertexID, numOfPartitions uint64) (uint64, error) {
			return plugin.HashPartition(id, numOfPartitions)
		},
		GetAggregatorsMock: func() []plugin.Aggregator {
			return []plugin.Aggregator{aggregator.VertexStatsAggregatorInstance}
		},
	}

	workerProps := actor.PropsFromFunc(func(c actor.Context) {
		mux.Lock()
		defer mux.Unlock()
		switch cmd := c.Message().(type) {
		case *command.InitWorker:
			initCount++
			receivedPartitions = append(receivedPartitions, cmd.Partitions...)
			c.Respond(&command.InitWorkerAck{WorkerPid: c.Self()})
			if initCount == 3 {
				waitCh <- "InitWorker"
				initCount = 0
			}
		case *command.LoadVertex:
			loadCount++
			c.Respond(&command.LoadVertexAck{VertexId: cmd.VertexId})
		case *command.SuperStepBarrier:
			barrierCount++
			if len(cmd.ClusterInfo.WorkerInfo) != 3 {
				t.Fatal("unexpected worker len")
			}
			c.Respond(&command.SuperStepBarrierWorkerAck{WorkerPid: c.Self()})
			if barrierCount == 3 {
				waitCh <- "SuperStepBarrier"
				barrierCount = 0
			}
		case *command.Compute:
			stepCount++
			active := 2
			if cmd.SuperStep == 2 {
				active = 0 // to finish
			}
			v, err := aggregator.VertexStatsAggregatorInstance.MarshalValue(&aggregator.VertexStats{
				ActiveVertices: uint64(active),
				TotalVertices:  2,
				MessagesSent:   uint64(active * 2),
			})
			if err != nil {
				t.Fatal(err)
			}
			c.Respond(&command.ComputeWorkerAck{
				WorkerPid: c.Self(),
				AggregatedValues: map[string]*types.Any{
					aggregator.VertexStatsName: v,
				},
			})
			if stepCount == 3 {
				waitCh <- "Compute"
				stepCount = 0
			}
		}
	})

	coordinatorProps := actor.PropsFromProducer(func() actor.Actor {
		return NewCoordinatorActor(plugin, workerProps, logger)
	})
	context := actor.EmptyRootContext
	proxy := util.NewActorProxy(context, coordinatorProps, func(ctx actor.Context) {
	})

	// initialize
	initCount = 0
	proxy.Send(context, &command.NewCluster{
		Workers: []*command.NewCluster_WorkerReq{
			{Remote: false},
			{Remote: false},
			{Remote: false},
		},
		NrOfPartitions: 10,
	})

	println("wait for InitWorker")
	if s := <-waitCh; s != "InitWorker" {
		t.Fatal("unexpected initCount")
	}
	sort.Slice(receivedPartitions, func(i, j int) bool { return receivedPartitions[i] < receivedPartitions[j] })
	if diff := cmp.Diff(receivedPartitions, []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}); diff != "" {
		t.Fatalf("unexpected partitions: %s", diff)
	}

	if _, err := proxy.SendAndAwait(context, &command.LoadVertex{VertexId: "a"}, &command.LoadVertexAck{}, time.Second); err != nil {
		t.Fatal(err)
	}
	if _, err := proxy.SendAndAwait(context, &command.LoadVertex{VertexId: "b"}, &command.LoadVertexAck{}, time.Second); err != nil {
		t.Fatal(err)
	}
	if loadCount != 2 {
		t.Fatal("unexpected load count")
	}

	// start
	proxy.Send(context, &command.StartSuperStep{})

	// step 0
	println("wait for SuperStepBarrier 0")
	if s := <-waitCh; s != "SuperStepBarrier" {
		t.Fatal("unexpected barrierCount")
	}
	println("wait for Compute 0")
	if s := <-waitCh; s != "Compute" {
		t.Fatal("unexpected stepCount")
	}

	resp, err := proxy.SendAndAwait(context, &command.CoordinatorStats{}, &command.CoordinatorStatsAck{}, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(resp, &command.CoordinatorStatsAck{
		SuperStep:        1,
		NrOfActiveVertex: 6,
		NrOfSentMessages: 12,
		State:            "processing superstep",
	}); diff != "" {
		t.Fatalf("unexpected stats: %s", diff)
	}

	// step 1
	println("wait for SuperStepBarrier 1")
	if s := <-waitCh; s != "SuperStepBarrier" {
		t.Fatal("unexpected barrierCount")
	}
	println("wait for Compute 1")
	if s := <-waitCh; s != "Compute" {
		t.Fatal("unexpected stepCount")
	}

	// step 2
	println("wait for SuperStepBarrier 2")
	if s := <-waitCh; s != "SuperStepBarrier" {
		t.Fatal("unexpected barrierCount")
	}
	println("wait for Compute 2")
	if s := <-waitCh; s != "Compute" {
		t.Fatal("unexpected stepCount")
	}

	resp, err = proxy.SendAndAwait(context, &command.CoordinatorStats{}, &command.CoordinatorStatsAck{}, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(resp, &command.CoordinatorStatsAck{
		SuperStep:        2,
		NrOfActiveVertex: 0,
		NrOfSentMessages: 0,
		State:            "idle",
	}); diff != "" {
		t.Fatalf("unexpected stats: %s", diff)
	}
}
