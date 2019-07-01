package main

import (
	"context"
	"fmt"

	"github.com/rerorero/prerogel/examples/sssp/loader"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/worker"
)

func main() {
	var plg plugin.Plugin

	plg = &ssspPlugin{
		sourceID: "a", // TODO: how should I specify the source vertex
		graph:    &loader.HeapLoader{},
	}

	fmt.Println("start agent..")
	if err := worker.Run(context.Background(), plg, ""); err != nil {
		panic(err)
	}
}
