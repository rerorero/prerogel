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
	conf, err := config.LoadWorkerConfFromEnv(envPrefix)
	if err != nil {
		return err
	}

	switch c := conf.(type) {
	case *config.MasterEnv:
		return RunMaster(ctx, plg, c)
	case *config.WorkerEnv:
		return RunWorker(ctx, plg, c)
	default:
		return fmt.Errorf("invalid config: %#v", c)
	}
}

// RunMaster starts running as master worker
func RunMaster(ctx context.Context, plg plugin.Plugin, conf *config.MasterEnv) error {
	root := actor.EmptyRootContext
	logger := conf.Logger()
	wait := newWaiting(ctx)

	// injection aggregators used for internal
	plg = newPluginProxy(plg).appendAggregators(systemAggregator)

	workerForLocal := workerProps(plg, logger, nil)
	coordinatorProps := actor.PropsFromProducer(func() actor.Actor {
		return NewCoordinatorActor(plg, workerForLocal, wait.shutdownHandler, logger)
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

	logger.Info(fmt.Sprintf("coordinator is running: addr=%s partitions=%v workers=%s log=%s",
		conf.ListenAddress, conf.Partitions, conf.WorkerAddresses, logger.Level.String()))

	wait.waitUntilDone()
	return nil
}

// RunWorker starts running as normal worker
func RunWorker(ctx context.Context, plg plugin.Plugin, conf *config.WorkerEnv) error {
	logger := conf.Logger()
	wait := newWaiting(ctx)
	// injection aggregators used for internal
	plg = newPluginProxy(plg).appendAggregators(systemAggregator)

	remote.Register(WorkerActorKind, workerProps(plg, logger, wait))
	remote.Start(conf.ListenAddress)

	logger.Info(fmt.Sprintf("worker is running: addr=%s log=%s", conf.ListenAddress, logger.Level.String()))

	wait.waitUntilDone()
	return nil
}

func workerProps(plg plugin.Plugin, logger *logrus.Logger, w *waiting) *actor.Props {
	vertexProps := actor.PropsFromProducer(func() actor.Actor {
		return NewVertexActor(plg, logger)
	})
	partitionProps := actor.PropsFromProducer(func() actor.Actor {
		return NewPartitionActor(plg, vertexProps, logger)
	})
	return actor.PropsFromProducer(func() actor.Actor {
		return NewWorkerActor(plg, partitionProps, w.shutdownHandler, logger)
	})
}

type waiting struct {
	shutdownCh chan struct{}
	ctx        context.Context
}

func newWaiting(ctx context.Context) *waiting {
	return &waiting{
		ctx:        ctx,
		shutdownCh: make(chan struct{}, 1),
	}
}

func (w *waiting) shutdownHandler() {
	if w != nil {
		w.shutdownCh <- struct{}{}
	}
}

func (w *waiting) waitUntilDone() {
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	defer close(w.shutdownCh)

	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		logrus.Warn("received SIGTERM")
	case <-w.shutdownCh:
		time.Sleep(1 * time.Second) // wait for message has sent
	case <-w.ctx.Done():
	}
}
