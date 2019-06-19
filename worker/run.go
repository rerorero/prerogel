package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/rerorero/prerogel/plugin"
	"github.com/sirupsen/logrus"
)

// Run starts worker server
func Run(ctx context.Context, plg plugin.Plugin, envPrefix string) error {
	conf, err := ReadEnv(envPrefix)
	if err != nil {
		return err
	}

	switch c := conf.(type) {
	case *MasterEnv:
		root := actor.EmptyRootContext
		logger := c.Logger()

		// init worker

		remote.Start(c.ListenAddress)

	case *WorkerEnv:
		logger := c.Logger()

		vertexProps := actor.PropsFromProducer(func() actor.Actor {
			return NewVertexActor(plg, logger)
		})
		partitionProps := actor.PropsFromProducer(func() actor.Actor {
			return NewPartitionActor(plg, vertexProps, logger)
		})
		workerProps := actor.PropsFromProducer(func() actor.Actor {
			return NewWorkerActor(plg, partitionProps, logger)
		})

		remote.Register(RemoteWorkerName, workerProps)
		remote.Start(c.ListenAddress)

	default:
		return fmt.Errorf("invalid config: %#v", c)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		logrus.Warn("received SIGTERM")
	case <-ctx.Done():
	}

	return nil
}
