package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/config"
	"github.com/rerorero/prerogel/plugin"
	"github.com/sirupsen/logrus"
)

// Run starts worker server
func Run(ctx context.Context, plg plugin.Plugin, envPrefix string) error {
	println("natoring run")
	conf, err := config.LoadWorkerConfFromEnv(envPrefix)
	if err != nil {
		return err
	}

	switch c := conf.(type) {
	case *config.MasterEnv:
		println("natoring mas")
		return RunMaster(ctx, plg, c)
	case *config.WorkerEnv:
		println("natoring wor")
		return RunWorker(ctx, plg, c)
	default:
		return fmt.Errorf("invalid config: %#v", c)
	}
}

// RunMaster starts running as master worker
func RunMaster(ctx context.Context, plg plugin.Plugin, conf *config.MasterEnv) error {
	root := actor.EmptyRootContext
	logger := conf.Logger()

	workerForLocal := workerProps(plg, logger)
	coordinatorProps := actor.PropsFromProducer(func() actor.Actor {
		return NewCoordinatorActor(plg, workerForLocal, logger)
	})

	remote.Start(conf.ListenAddress)

	coordinator, err := root.SpawnNamed(coordinatorProps, CoordinatorActorID)
	if err != nil {
		return errors.Wrapf(err, "failed to spawn coordinator actor")
	}

	var workers []*command.NewCluster_WorkerReq
	for _, w := range conf.WorkerAddresses {
		workers = append(workers, &command.NewCluster_WorkerReq{
			Remote:      true,
			HostAndPort: w,
		})
	}

	f := root.RequestFuture(coordinator, &command.NewCluster{
		Workers:        workers,
		NrOfPartitions: conf.Partitions,
	}, 120*time.Second)
	if err := f.Wait(); err != nil {
		return errors.Wrap(err, "failed to connect to worker: ")
	}
	res, err := f.Result()
	if err != nil {
		return errors.Wrap(err, "failed to initialize coordinator: ")
	}
	if _, ok := res.(*command.NewClusterAck); !ok {
		return fmt.Errorf("failed to initialize cooridnator: unknown ack %#v", res)
	}

	logger.Info(fmt.Sprintf("coordinator is running: addr=%s log=%s", conf.ListenAddress, logger.Level.String()))

	waitUntilDone(ctx)
	return nil
}

// RunWorker starts running as normal worker
func RunWorker(ctx context.Context, plg plugin.Plugin, conf *config.WorkerEnv) error {
	println("natoring1")
	logger := conf.Logger()
	remote.Register(WorkerActorKind, workerProps(plg, logger))
	println("natoring2")
	remote.Start(conf.ListenAddress)
	println("natoring3")

	logger.Info(fmt.Sprintf("worker is running: addr=%s log=%s", conf.ListenAddress, logger.Level.String()))

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
