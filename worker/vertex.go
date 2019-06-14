package worker

import (
	"fmt"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/util"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus"
)

type vertexActor struct {
	util.ActorUtil
	behavior              actor.Behavior
	plugin                Plugin
	vertex                Vertex
	partitionID           uint64
	halted                bool
	prevStepMessages      []Message
	messageQueue          []Message
	ackRecorder           *util.AckRecorder
	computeRespondTo      *actor.PID
	aggregatedCurrentStep []*command.AggregatedValue
}

type computeContextImpl struct {
	superStep          uint64
	ctx                actor.Context
	vertexActor        *vertexActor
	aggregatedPrevStep []*command.AggregatedValue
}

func (c *computeContextImpl) SuperStep() uint64 {
	return c.superStep
}

func (c *computeContextImpl) ReceivedMessages() []Message {
	return c.vertexActor.prevStepMessages
}

func (c *computeContextImpl) SendMessageTo(dest VertexID, m Message) error {
	pb, err := c.vertexActor.plugin.MarshalMessage(m)
	if err != nil {
		c.vertexActor.ActorUtil.LogError(fmt.Sprintf("failed to marshal message: id=%v, message=%#v", c.vertexActor.vertex.GetID(), m))
		return err
	}
	messageID := uuid.New().String()
	c.ctx.Request(c.ctx.Parent(), &command.SuperStepMessage{
		Uuid:         messageID,
		SuperStep:    c.superStep,
		SrcVertexId:  string(c.vertexActor.vertex.GetID()),
		DestVertexId: string(dest),
		Message:      pb,
	})

	if !c.vertexActor.ackRecorder.AddToWaitList(messageID) {
		c.vertexActor.ActorUtil.LogWarn(fmt.Sprintf("duplicate superstep message: from=%v to=%v", c.vertexActor.vertex.GetID(), dest))
	}

	return nil
}

func (c *computeContextImpl) VoteToHalt() {
	c.vertexActor.ActorUtil.LogDebug("halted by user")
	c.vertexActor.halted = true
}

func (c *computeContextImpl) GetAggregated(aggregatorName string) (AggregatableValue, bool, error) {
	aggregator, err := c.findAggregator(aggregatorName)
	if err != nil {
		return nil, false, err
	}
	for _, a := range c.aggregatedPrevStep {
		if a.AggregatorName == aggregatorName {
			v, err := aggregator.UnmarshalValue(a.Value)
			if err != nil {
				return nil, false, errors.Wrapf(err, "failed to unmarshal aggregated value: %+v", a.Value)
			}
			return v, true, nil
		}
	}
	return nil, false, nil
}

func (c *computeContextImpl) PutAggregatable(aggregatorName string, v AggregatableValue) error {
	aggregator, err := c.findAggregator(aggregatorName)
	if err != nil {
		return err
	}

	for _, current := range c.vertexActor.aggregatedCurrentStep {
		if current.AggregatorName == aggregatorName {
			// aggregate TODO: verbose marshalling
			v2, err := aggregator.UnmarshalValue(current.Value)
			if err != nil {
				return errors.Wrapf(err, "failed to unmarshal aggregated value: %+v", current.Value)
			}
			pb, err := aggregator.MarshalValue(aggregator.Aggregate(v, v2))
			if err != nil {
				return errors.Wrapf(err, "failed to marshal aggregatable value: %#v", v)
			}
			current.Value = pb
			return nil
		}
	}

	pb, err := aggregator.MarshalValue(v)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal aggregatable value: %#v", v)
	}
	c.vertexActor.aggregatedCurrentStep = append(c.vertexActor.aggregatedCurrentStep, &command.AggregatedValue{
		AggregatorName: aggregatorName,
		Value:          pb,
	})
	return nil
}

func (c *computeContextImpl) findAggregator(name string) (Aggregator, error) {
	for _, a := range c.vertexActor.plugin.GetAggregators() {
		if a.Name() == name {
			return a, nil
		}
	}
	return nil, fmt.Errorf("%s: no such aggregator", name)
}

// NewVertexActor returns an actor instance
func NewVertexActor(plugin Plugin, logger *logrus.Logger) actor.Actor {
	ar := &util.AckRecorder{}
	ar.Clear()
	a := &vertexActor{
		plugin: plugin,
		ActorUtil: util.ActorUtil{
			Logger: logger,
		},
		ackRecorder: ar,
	}
	a.behavior.Become(a.waitInit)
	return a
}

// Receive is message handler
func (state *vertexActor) Receive(context actor.Context) {
	if state.ActorUtil.IsSystemMessage(context.Message()) {
		return
	}
	state.behavior.Receive(context)
}

func (state *vertexActor) waitInit(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.LoadVertex:
		if state.vertex != nil {
			state.ActorUtil.Fail(errors.New("vertex has already initialized"))
			return
		}
		state.vertex = state.plugin.NewVertex(VertexID(cmd.VertexId))
		if err := state.vertex.Load(); err != nil {
			state.ActorUtil.Fail(errors.New("failed to load vertex"))
			return
		}
		state.partitionID = cmd.PartitionId
		state.ActorUtil.AppendLoggerField("vertexId", state.vertex.GetID())
		context.Respond(&command.LoadVertexAck{
			PartitionId: cmd.PartitionId,
			VertexId:    string(state.vertex.GetID()),
		})
		state.behavior.Become(state.superstep)
		state.ActorUtil.LogDebug("vertex has initialized")
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[waitInit] unhandled vertex command: command=%#v(%v)", cmd, reflect.TypeOf(cmd)))
		return
	}
}

func (state *vertexActor) superstep(context actor.Context) {
	switch cmd := context.Message().(type) {
	case *command.SuperStepBarrier:
		// move messages from queue to buffer
		state.prevStepMessages = state.messageQueue
		state.messageQueue = nil
		context.Respond(&command.SuperStepBarrierAck{
			VertexId: string(state.vertex.GetID()),
		})
		state.ActorUtil.LogDebug(fmt.Sprintf("received barrier message"))
		return

	case *command.Compute:
		state.computeRespondTo = context.Sender()
		state.onComputed(context, cmd)
		return

	case *command.SuperStepMessage:
		if state.vertex.GetID() != VertexID(cmd.DestVertexId) {
			state.ActorUtil.Fail(fmt.Errorf("inconsistent vertex id: %#v", *cmd))
			return
		}
		// TODO: verify cmd.SuperStep
		pb, err := state.plugin.UnmarshalMessage(cmd.Message)
		if err != nil {
			state.ActorUtil.Fail(fmt.Errorf("failed to unmarshal message: %#v", *cmd))
			return
		}
		state.messageQueue = append(state.messageQueue, pb)
		state.halted = false
		context.Respond(&command.SuperStepMessageAck{
			Uuid: cmd.Uuid,
		})
		return

	case *command.SuperStepMessageAck:
		if state.ackRecorder.HasCompleted() {
			state.ActorUtil.LogWarn(fmt.Sprintf("unhaneled message id=%v, compute() has already completed", cmd.Uuid))
			return
		}
		if !state.ackRecorder.Ack(cmd.Uuid) {
			state.ActorUtil.LogWarn(fmt.Sprintf("duplicated or unhaneled message: uuid=%v", cmd.Uuid))
		}
		if state.ackRecorder.HasCompleted() {
			state.respondComputeAck(context)
		}
		return

	default:
		state.ActorUtil.Fail(fmt.Errorf("[superstep] unhandled vertex command: command=%#v", cmd))
		return
	}
}

func (state *vertexActor) onComputed(ctx actor.Context, cmd *command.Compute) {
	// force to compute() in super step 0
	// otherwise halt if there are no messages
	if cmd.SuperStep == 0 {
		state.halted = false
	} else if len(state.prevStepMessages) == 0 {
		state.ActorUtil.LogDebug("halted due to no message")
		state.halted = true
	}

	if state.halted {
		state.respondComputeAck(ctx)
		return
	}

	state.ackRecorder.Clear()
	computeContext := &computeContextImpl{
		superStep:          cmd.SuperStep,
		ctx:                ctx,
		vertexActor:        state,
		aggregatedPrevStep: cmd.AggregatedValues,
	}
	if err := state.vertex.Compute(computeContext); err != nil {
		state.ActorUtil.Fail(errors.Wrap(err, "failed to compute"))
		return
	}

	if state.ackRecorder.HasCompleted() {
		state.respondComputeAck(ctx)
	}
	return
}

func (state *vertexActor) respondComputeAck(ctx actor.Context) {
	ctx.Send(state.computeRespondTo, &command.ComputeAck{
		VertexId:         string(state.vertex.GetID()),
		Halted:           state.halted,
		AggregatedValues: state.aggregatedCurrentStep,
	})
	state.ActorUtil.LogDebug("compute() completed")
}
