package main

import (
	"context"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/config"
	"github.com/rerorero/prerogel/examples/maximum/loader"
	"github.com/rerorero/prerogel/worker"
)

func Test_maximumPlugin(t *testing.T) {
	graph := &loader.HeapLoader{
		Vertices: map[string]uint32{
			"a": 3,
			"b": 6,
			"c": 2,
			"d": 1,
		},
		Edges: [][]string{
			// src, dest
			{"a", "b"},
			{"b", "a"},
			{"b", "d"},
			{"c", "b"},
			{"c", "d"},
			{"d", "c"},
		},
	}

	plg := &maxPlugin{
		graph: graph,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	masterAddr := "127.0.0.1:8990"
	workerAddrress := []string{"127.0.0.1:8991", "127.0.0.1:8992"}

	go func() {
		for _, addr := range workerAddrress {
			if err := worker.RunWorker(ctx, plg, &config.WorkerEnv{
				CommonConfig:  config.CommonConfig{LogLevel: "DEBUG"},
				ListenAddress: addr,
			}); err != nil {
				t.Fatal(err)
			}
		}
	}()

	go func() {
		if err := worker.RunMaster(ctx, plg, &config.MasterEnv{
			CommonConfig:    config.CommonConfig{LogLevel: "DEBUG"},
			ListenAddress:   masterAddr,
			WorkerAddresses: workerAddrress,
			Partitions:      3,
		}); err != nil {
			t.Fatal(err)
		}
	}()

	c := actor.EmptyRootContext
	pid := &actor.PID{
		Address: masterAddr,
		Id:      worker.CoordinatorActorID,
	}

	ticker := time.NewTicker(500 * time.Millisecond)

LOOP:
	for {
		select {
		case <-ctx.Done():
			t.Fatal("timeout")
		case <-ticker.C:
			cmd := &command.CoordinatorStats{}

			fut := c.RequestFuture(pid, cmd, 5*time.Second)
			if err := fut.Wait(); err != nil {
				t.Fatal(err)
			}
			res, err := fut.Result()
			if err != nil {
				t.Fatal(err)
			}
			stat, ok := res.(*command.CoordinatorStatsAck)
			if !ok {
				t.Fatalf("invalid CoordinatorStatsAck: %#v", res)
			}

			if stat.StatsCompleted() {
				break LOOP
			}
		}
	}
}
