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

	plg = loader.HeapLoader{}

	fmt.Println("start agent..")
	if err := worker.Run(context.Background(), plg, ""); err != nil {
		panic(err)
	}
}
