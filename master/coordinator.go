package master

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus"
)

type coordinatorActor struct {
	util.ActorUtil
	behavior              actor.Behavior
	plugin                plugin.Plugin
	workerProps           *actor.Props
	clusterInfo           *command.ClusterInfo
	ackRecorder           *util.AckRecorder
	aggregatedCurrentStep map[string]*any.Any
}

// NewCoordinatorActor returns an actor instance
func NewCoordinatorActor(plg plugin.Plugin, workerProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &util.AckRecorder{}
	ar.Clear()
	a := &coordinatorActor{
		plugin: plg,
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		workerProps: workerProps,
		ackRecorder: ar,
	}
	a.behavior.Become(a.waitInit)
	return a
}

// Receive is message handler
func (state *coordinatorActor) Receive(context actor.Context) {
	if state.ActorUtil.IsSystemMessage(context.Message()) {
		// ignore
		return
	}
	state.behavior.Receive(context)
}

func (state *coordinatorActor) waitInit(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.AddWorker:

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled worker command: command=%#v", cmd))
		return
	}
}
