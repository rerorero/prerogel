package worker

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rerorero/prerogel/util"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus"
)

type workerActor struct {
	util.ActorUtil
	behavior       actor.Behavior
	plugin         Plugin
	partitions     map[uint64]*actor.PID
	partitionProps *actor.Props
	clusterInfo    *command.ClusterInfo
	ackRecorder    *ackRecorder
}

// NewWorkerActor returns a new actor instance
func NewWorkerActor(plugin Plugin, partitionProps *actor.Props, logger *logrus.Logger) actor.Actor {
	ar := &ackRecorder{}
	ar.clear()
	a := &workerActor{
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		plugin:         plugin,
		partitions:     make(map[uint64]*actor.PID),
		partitionProps: partitionProps,
		ackRecorder:    ar,
	}
	a.behavior.Become(a.waitInit)
	return a
}

// Receive is message handler
func (state *workerActor) Receive(context actor.Context) {
	if state.ActorUtil.IsSystemMessage(context.Message()) {
		// ignore
		return
	}
	state.behavior.Receive(context)
}

func (state *workerActor) waitInit(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitWorker:
		contains := false
		for _, pm := range cmd.ClusterInfo.PartitionMap {
			if pm.WorkerPid.Id == context.Self().GetId() {
				contains = true
				break
			}
		}
		if !contains {
			state.ActorUtil.Fail(fmt.Errorf("worker not assigned: me=%v", context.Self().GetId()))
			return
		}
		state.ActorUtil.AppendLoggerField("worker_id", context.Self().GetId())
		state.clusterInfo = cmd.ClusterInfo
		context.Respond(&command.InitWorkerAck{
			WorkerPid: context.Self(),
		})
		state.behavior.Become(state.run)
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled worker command: command=%+v", cmd))
		return
	}
}

func (state *workerActor) run(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.InitPartition:
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled worker command: command=%+v", cmd))
		return
	}
}
