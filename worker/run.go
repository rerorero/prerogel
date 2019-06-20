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
		return RunMaster(ctx, plg, c)
	case *WorkerEnv:
		return RunWorker(ctx, plg, c)
	default:
		return fmt.Errorf("invalid config: %#v", c)
	}

	return nil
}

// RunMaster starts running as master worker
func RunMaster(ctx context.Context, plg plugin.Plugin, conf *MasterEnv) error {
	root := actor.EmptyRootContext
	logger := conf.Logger()

	workerForLocal := workerProps(plg, logger)
	coordinator := actor.PropsFromProducer(func() actor.Actor {
		return NewCoordinatorActor(plg, workerForLocal, logger)
	})

	remote.Start(conf.ListenAddress)

	root.SpawnNamed(coordinator, CoordinatorActorKind)

	waitUntilDone(ctx)
	return nil
}

// RunWorker starts running as normal worker
func RunWorker(ctx context.Context, plg plugin.Plugin, conf *WorkerEnv) error {
	remote.Register(WorkerActorKind, workerProps(plg, conf.Logger()))
	remote.Start(conf.ListenAddress)

	waitUntilDone(ctx)
	return nil
}

func workerProps(plg plugin.Plugin, logger *logrus.Logger) *actor.Props {
	vertexProps := actor.PropsFromProducer(func() actor.Actor {
		return NewVertexActor(plg, logger)
	})
	partitionProps := actor.PropsFromProducer(func() actor.Actor {
		return NewPartitionActor(plg, vertexProps, logger)
	})
	return actor.PropsFromProducer(func() actor.Actor {
		return NewWorkerActor(plg, partitionProps, logger)
	})
}

func waitUntilDone(ctx context.Context) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		logrus.Warn("received SIGTERM")
	case <-ctx.Done():
	}
}
