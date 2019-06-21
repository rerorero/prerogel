package maximum

import (
	"context"
	"testing"
	"time"

	"github.com/rerorero/prerogel/config"

	"github.com/rerorero/prerogel/worker"

	"github.com/rerorero/prerogel/examples/plugins/maximum/loader"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	workerAddrress := []string{"127.0.0.1:8990", "127.0.0.1:8991"}

	for _, addr := range workerAddrress {
		go func() {

		}
		if err := worker.RunWorker(ctx, plg, &config.WorkerEnv{
			CommonConfig:  config.CommonConfig{LogLevel: "DEBUG"},
			ListenAddress: addr,
		}); err != nil {
			t.Fatal(err)
		}
	}

}
