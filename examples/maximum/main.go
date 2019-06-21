package main

import (
	"context"
	"fmt"

	"github.com/rerorero/prerogel/examples/maximum/loader"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/worker"
)

func main() {
	var plg plugin.Plugin

	plg = pluginWithHeap()

	fmt.Println("start agent..")
	if err := worker.Run(context.Background(), plg, ""); err != nil {
		panic(err)
	}
}

func pluginWithHeap() plugin.Plugin {
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

	return &maxPlugin{
		graph: graph,
	}
}
