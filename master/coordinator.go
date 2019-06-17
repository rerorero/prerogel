package master

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
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
	a.behavior.Become(a.idle)
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

func (state *coordinatorActor) idle(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.NewCluster:
		//for _, : = range cmd.Workers

		//if cmd.HostAndPort == "" {
		//	pid = context.Spawn(state.workerProps)
		//} else {
		//	pidRes, err := remote.SpawnNamed(cmd.HostAndPort, "master", "worker", 30*time.Second)
		//	if err != nil {
		//		state.ActorUtil.LogError(fmt.Sprintf("failed to spawn remote actor: %d", pidRes.StatusCode))
		//	}
		//	if err != nil {
		//		state.ActorUtil.LogError(fmt.Sprintf("failed to spawn remote actor: %d", pidRes.StatusCode))
		//		return
		//	}
		//	pid = pidRes.Pid
		//}
	default:
		state.ActorUtil.Fail(fmt.Errorf("[idle] unhandled corrdinator command: command=%#v", cmd))
		return
	}
}

type workerAndPartitions struct {
	worker     string
	partitions []uint64
}

func assignPartition(workers []string, partitions uint64) ([]*workerAndPartitions, error) {
	if len(workers) == 0 {
		return nil, errors.New("no available workers")
	}
	var pairs []*workerAndPartitions
	max := partitions / uint64(len(workers))

	return pairs
}
