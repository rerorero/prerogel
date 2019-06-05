package worker

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus"
)

type partitionActor struct {
	behavior actor.Behavior
	plugin   Plugin
	logger   *logrus.Logger
	vertices map[VertexID]*actor.PID
}

// NewVertexActor returns an actor instance
func NewPartitionActor(plugin Plugin, logger *logrus.Logger) actor.Actor {
	a := &partitionActor{
		plugin: plugin,
		logger: logger,
	}
	a.behavior.Become(a.waitInit)
	return a
}

// Receive is message handler
func (state *partitionActor) Receive(context actor.Context) {
	state.behavior.Receive(context)
}

func (state *partitionActor) waitInit(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitPartition:
		if ids, err := state.plugin.ListVertexID(cmd.PartitionId); err != nil {

		}
	default:
	}
}

func (state *partitionActor) fail(err error) {
	state.logger.WithError(err).Error(err.Error())
	// let it crash
	panic(err)
}

func (state *partitionActor) logError(msg string) {
	state.logger.Error(msg)
}

func (state *partitionActor) logInfo(msg string) {
	state.logger.Info(msg)
}
