package worker

import (
	"fmt"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/aggregator"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus"
)

type vertexActor struct {
	util.ActorUtil
	behavior              actor.Behavior
	plugin                plugin.Plugin
	vertex                plugin.Vertex
	partitionID           uint64
	halted                bool
	prevStepMessages      []plugin.Message
	messageQueue          []plugin.Message
	ackRecorder           *util.AckRecorder
	computeRespondTo      *actor.PID
	aggregatedCurrentStep map[string]*types.Any
	statsMessageSent      uint64
}

type loadVertexLocal struct {
	vertex plugin.Vertex
}

type computeContextImpl struct {
	superStep          uint64
	ctx                actor.Context
	vertexActor        *vertexActor
	aggregatedPrevStep map[string]*types.Any
}

var _ = (plugin.ComputeContext)(&computeContextImpl{})

func (c *computeContextImpl) SuperStep() uint64 {
	return c.superStep
}

func (c *computeContextImpl) ReceivedMessages() []plugin.Message {
	return c.vertexActor.prevStepMessages
}

func (c *computeContextImpl) SendMessageTo(dest plugin.VertexID, m plugin.Message) error {
	pb, err := c.vertexActor.plugin.MarshalMessage(m)
	if err != nil {
		c.vertexActor.ActorUtil.LogError(c.ctx, fmt.Sprintf("failed to marshal message: id=%v, message=%#v", c.vertexActor.vertex.GetID(), m))
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

	c.vertexActor.ActorUtil.LogDebug(c.ctx, fmt.Sprintf("message sent: uuid=%s, %v -> %v",
		messageID, c.vertexActor.vertex.GetID(), dest))

	if !c.vertexActor.ackRecorder.AddToWaitList(messageID) {
		c.vertexActor.ActorUtil.LogWarn(c.ctx, fmt.Sprintf("duplicate superstep message: from=%v to=%v", c.vertexActor.vertex.GetID(), dest))
	}

	return nil
}

func (c *computeContextImpl) VoteToHalt() {
	c.vertexActor.ActorUtil.LogDebug(c.ctx, "halted by user")
	c.vertexActor.halted = true
}

func (c *computeContextImpl) GetAggregated(aggregatorName string) (plugin.AggregatableValue, bool, error) {
	aggregator, err := findAggregator(c.vertexActor.plugin.GetAggregators(), aggregatorName)
	if err != nil {
		return nil, false, err
	}

	value, ok := c.aggregatedPrevStep[aggregatorName]
	if ok {
		v, err := aggregator.UnmarshalValue(value)
		if err != nil {
			return nil, false, errors.Wrapf(err, "failed to unmarshal aggregated value: %+v", value)
		}
		return v, true, nil
	}
	return nil, false, nil
}

func (c *computeContextImpl) PutAggregatable(aggregatorName string, v plugin.AggregatableValue) error {
	aggregator, err := findAggregator(c.vertexActor.plugin.GetAggregators(), aggregatorName)
	if err != nil {
		return err
	}

	current, ok := c.vertexActor.aggregatedCurrentStep[aggregatorName]
	if ok {
		// aggregate TODO: verbose marshalling
		v2, err := aggregator.UnmarshalValue(current)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal aggregated value: %+v", current.Value)
		}
		val, err := aggregator.Aggregate(v, v2)
		if err != nil {
			return errors.Wrap(err, "failed to Aggregate()")
		}
		pb, err := aggregator.MarshalValue(val)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal aggregatable value: %#v", v)
		}
		c.vertexActor.aggregatedCurrentStep[aggregatorName] = pb
		return nil
	}

	pb, err := aggregator.MarshalValue(v)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal aggregatable value: %#v", v)
	}
	c.vertexActor.aggregatedCurrentStep[aggregatorName] = pb
	return nil
}

// NewVertexActor returns an actor instance
func NewVertexActor(plugin plugin.Plugin, logger *logrus.Logger) actor.Actor {
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

	switch cmd := context.Message().(type) {
	case *command.GetVertexValue:
		ack := &command.GetVertexValueAck{
			VertexId: cmd.VertexId,
		}
		if state.vertex != nil {
			ack.Value = state.vertex.GetValueAsString()
		}
		context.Respond(ack)
		return

	default:
		state.behavior.Receive(context)
		return
	}
}

func (state *vertexActor) waitInit(context actor.Context) {
	var vert plugin.Vertex

	switch cmd := context.Message().(type) {
	case *command.LoadVertex:
		var err error
		vert, err = state.plugin.NewVertex(plugin.VertexID(cmd.VertexId))
		if err != nil {
			e := fmt.Sprintf("failed to NewVertex: id=%s err=%s", cmd.VertexId, err.Error())
			state.ActorUtil.LogError(context, e)
			context.Respond(&command.LoadVertexAck{VertexId: string(state.vertex.GetID()), Error: e})
			return
		}

	case *loadVertexLocal:
		vert = cmd.vertex

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[waitInit] unhandled vertex command: command=%#v(%v)", cmd, reflect.TypeOf(cmd)))
	}

	if state.vertex != nil {
		err := fmt.Sprintf("vertex has already initialized: id=%s", vert.GetID())
		state.ActorUtil.LogError(context, err)
		context.Respond(&command.LoadVertexAck{VertexId: string(state.vertex.GetID()), Error: err})
		return
	}

	state.vertex = vert
	context.Respond(&command.LoadVertexAck{
		VertexId: string(state.vertex.GetID()),
	})
	state.behavior.Become(state.superstep)
	state.ActorUtil.LogInfo(context, fmt.Sprintf("vertex %v loaded", state.vertex.GetID()))
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
		state.ActorUtil.LogDebug(context, fmt.Sprintf("received barrier message"))
		return

	case *command.Compute:
		state.computeRespondTo = context.Sender()
		state.onComputed(context, cmd)
		return

	case *command.SuperStepMessage:
		if state.vertex.GetID() != plugin.VertexID(cmd.DestVertexId) {
			state.ActorUtil.Fail(context, fmt.Errorf("inconsistent vertex id: %#v", *cmd))
			return
		}
		// TODO: verify cmd.SuperStep
		pb, err := state.plugin.UnmarshalMessage(cmd.Message)
		if err != nil {
			state.ActorUtil.Fail(context, fmt.Errorf("failed to unmarshal message: %#v", *cmd))
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
			state.ActorUtil.LogWarn(context, fmt.Sprintf("unhaneled message id=%v, compute() has already completed", cmd.Uuid))
			return
		}
		if !state.ackRecorder.Ack(cmd.Uuid) {
			state.ActorUtil.LogWarn(context, fmt.Sprintf("duplicated or unhaneled message: uuid=%v", cmd.Uuid))
		}
		if state.ackRecorder.HasCompleted() {
			state.respondComputeAck(context)
		}
		return

	default:
		state.ActorUtil.Fail(context, fmt.Errorf("[superstep] unhandled vertex command: command=%#v", cmd))
		return
	}
}

func (state *vertexActor) onComputed(ctx actor.Context, cmd *command.Compute) {
	state.aggregatedCurrentStep = make(map[string]*types.Any)
	state.statsMessageSent = 0

	// force to compute() in super step 0
	// otherwise halt if there are no messages
	if cmd.SuperStep == 0 {
		state.halted = false
	} else if len(state.prevStepMessages) == 0 {
		state.ActorUtil.LogDebug(ctx, "halted due to no message")
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
		state.ActorUtil.Fail(ctx, errors.Wrap(err, "failed to compute"))
		return
	}
	state.statsMessageSent = uint64(state.ackRecorder.Size())

	if state.ackRecorder.HasCompleted() {
		state.respondComputeAck(ctx)
	}
	return
}

func (state *vertexActor) respondComputeAck(ctx actor.Context) {
	if len(state.messageQueue) > 0 {
		// activate if it receives messages to be handled in the next step
		state.halted = false
	}

	stats, err := findAggregator(state.plugin.GetAggregators(), VertexStatsName)
	if err == nil {
		var active uint64
		if !state.halted {
			active = 1
		}
		pb, err := stats.MarshalValue(&aggregator.VertexStats{
			ActiveVertices: active,
			TotalVertices:  1,
			MessagesSent:   state.statsMessageSent,
		})
		if err != nil {
			state.ActorUtil.LogError(ctx, fmt.Sprintf("failed to marshal stats: %v", err))
			return
		}
		state.aggregatedCurrentStep[VertexStatsName] = pb
	}

	// TODO: store vertex state

	ctx.Send(state.computeRespondTo, &command.ComputeAck{
		VertexId:         string(state.vertex.GetID()),
		Halted:           state.halted,
		AggregatedValues: state.aggregatedCurrentStep,
	})
	state.aggregatedCurrentStep = nil
	state.ActorUtil.LogDebug(ctx, "compute() completed")
}
