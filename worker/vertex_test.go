package worker

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus/hooks/test"
)

func Test_vertexActor_Receive_InitVertex(t *testing.T) {
	type fields struct {
		vertex plugin.Vertex
	}
	tests := []struct {
		name        string
		fields      fields
		cmd         []proto.Message
		wantRespond []proto.Message
		wantLoaded  int
	}{
		{
			name: "should be initialized",
			fields: fields{
				vertex: &MockedVertex{
					GetIDMock: func() plugin.VertexID {
						return "test-id"
					},
				},
			},
			cmd: []proto.Message{
				&command.LoadVertex{
					VertexId:    "test-id",
					PartitionId: 123,
				},
			},
			wantRespond: []proto.Message{
				&command.LoadVertexAck{
					VertexId:    "test-id",
					PartitionId: 123,
				},
			},
			wantLoaded: 1,
		},
		{
			name: "fails if already initialized",
			fields: fields{
				vertex: &MockedVertex{
					GetIDMock: func() plugin.VertexID {
						return "test-id"
					},
				},
			},
			cmd: []proto.Message{
				&command.LoadVertex{
					VertexId:    "test-id",
					PartitionId: 123,
				},
				&command.LoadVertex{},
			},
			wantRespond: []proto.Message{
				&command.LoadVertexAck{
					VertexId:    "test-id",
					PartitionId: 123,
				},
				nil,
			},
			wantLoaded: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := test.NewNullLogger()

			context := actor.EmptyRootContext
			props := actor.PropsFromProducer(func() actor.Actor {
				return NewVertexActor(&MockedPlugin{
					NewVertexMock: func(id plugin.VertexID) (plugin.Vertex, error) {
						if id != "test-id" {
							t.Fatal("unexpected vertex id")
						}
						return tt.fields.vertex, nil
					},
				}, logger)
			})
			pid := context.Spawn(props)

			for i, cmd := range tt.cmd {
				res, err := context.RequestFuture(pid, cmd, time.Second).Result()
				if (err != nil) != (tt.wantRespond[i] == nil) {
					t.Fatal(err)
				}
				if diff := cmp.Diff(tt.wantRespond[i], res); diff != "" {
					t.Errorf("unexpected respond: %s", diff)
				}
			}
		})
	}
}

func Test_vertexActor_Receive_Compute(t *testing.T) {
	var computed int
	msg1 := &command.SuperStepMessage{
		Uuid:         "a",
		SuperStep:    0,
		SrcVertexId:  "foo",
		DestVertexId: "test-id",
		Message:      &any.Any{TypeUrl: "com.example/test", Value: []byte("test1")},
	}
	msg2 := &command.SuperStepMessage{
		Uuid:         "b",
		SuperStep:    0,
		SrcVertexId:  "bar",
		DestVertexId: "test-id",
		Message:      &any.Any{TypeUrl: "com.example/test", Value: []byte("test2")},
	}

	tests := []struct {
		name             string
		vertex           plugin.Vertex
		cmd              []proto.Message
		incomingMessages map[int][]*command.SuperStepMessage
		wantRespond      []proto.Message
		wantComputed     int
		wantSentMessages []*command.SuperStepMessage
	}{
		{
			name: "compute 0 then halt",
			vertex: &MockedVertex{
				ComputeMock: func(ctx plugin.ComputeContext) error {
					if ctx.SuperStep() != 0 {
						t.Fatal("unexpected step")
					}
					if ctx.ReceivedMessages() != nil {
						t.Fatal("unexpected received messages")
					}
					computed++
					return nil
				},
				GetIDMock: func() plugin.VertexID {
					return "test-id"
				},
			},
			cmd: []proto.Message{
				&command.LoadVertex{VertexId: "test-id", PartitionId: 123},
				&command.SuperStepBarrier{},
				&command.Compute{SuperStep: 0},
				&command.Compute{SuperStep: 1},
			},
			wantRespond: []proto.Message{
				&command.LoadVertexAck{
					PartitionId: 123,
					VertexId:    "test-id",
				},
				&command.SuperStepBarrierAck{
					VertexId: string("test-id"),
				},
				&command.ComputeAck{
					VertexId:         string("test-id"),
					Halted:           false,
					AggregatedValues: make(map[string]*any.Any),
				},
				&command.ComputeAck{
					VertexId:         string("test-id"),
					Halted:           true,
					AggregatedValues: make(map[string]*any.Any),
				},
			},
			wantComputed:     1,
			wantSentMessages: nil,
		},
		{
			name: "receive messages and continue to compute",
			vertex: &MockedVertex{
				ComputeMock: func(ctx plugin.ComputeContext) error {
					if ctx.SuperStep() != uint64(computed) {
						t.Fatalf("unexpected step: %v, %v", ctx.SuperStep(), computed)
					}
					if ctx.SuperStep() == 1 {
						if diff := cmp.Diff([]plugin.Message{"test1", "test2"}, ctx.ReceivedMessages()); diff != "" {
							t.Fatalf("unexpected received messages: %s", diff)
						}
					} else if ctx.ReceivedMessages() != nil {
						t.Fatal("unexpected received messages")
					}
					computed++
					return nil
				},
				GetIDMock: func() plugin.VertexID { return "test-id" },
			},
			cmd: []proto.Message{
				&command.LoadVertex{VertexId: "test-id", PartitionId: 123},
				&command.SuperStepBarrier{},
				&command.Compute{SuperStep: 0},
				&command.SuperStepBarrier{},
				&command.Compute{SuperStep: 1},
				&command.SuperStepBarrier{},
				&command.Compute{SuperStep: 2},
			},
			incomingMessages: map[int][]*command.SuperStepMessage{
				2: {msg1, msg2}, // sent between Compute(step0) and Compute(step1)
			},
			wantRespond: []proto.Message{
				&command.LoadVertexAck{VertexId: "test-id", PartitionId: 123},
				&command.SuperStepBarrierAck{VertexId: string("test-id")},
				&command.ComputeAck{VertexId: string("test-id"), Halted: false, AggregatedValues: make(map[string]*any.Any)},
				&command.SuperStepBarrierAck{VertexId: string("test-id")},
				&command.ComputeAck{VertexId: string("test-id"), Halted: false, AggregatedValues: make(map[string]*any.Any)},
				&command.SuperStepBarrierAck{VertexId: string("test-id")},
				&command.ComputeAck{VertexId: string("test-id"), Halted: true, AggregatedValues: make(map[string]*any.Any)},
			},
			wantComputed:     2,
			wantSentMessages: nil,
		},
		{
			name: "send messages to other vertices",
			vertex: &MockedVertex{
				ComputeMock: func(ctx plugin.ComputeContext) error {
					msg := &types.Timestamp{
						Seconds: 123456,
						Nanos:   int32(ctx.SuperStep()),
					}
					if err := ctx.SendMessageTo(plugin.VertexID(fmt.Sprintf("dest-%v", ctx.SuperStep())), msg); err != nil {
						t.Fatal(err)
					}
					computed++
					return nil
				},
				GetIDMock: func() plugin.VertexID { return "test-id" },
			},
			cmd: []proto.Message{
				&command.LoadVertex{VertexId: "test-id", PartitionId: 123},
				&command.SuperStepBarrier{},
				&command.Compute{SuperStep: 0},
				&command.SuperStepBarrier{},
				&command.Compute{SuperStep: 1},
				&command.SuperStepBarrier{},
				&command.Compute{SuperStep: 2},
			},
			incomingMessages: map[int][]*command.SuperStepMessage{
				2: {msg1}, // sent between Compute(step0) and Compute(step1)
			},
			wantRespond: []proto.Message{
				&command.LoadVertexAck{VertexId: "test-id", PartitionId: 123},
				&command.SuperStepBarrierAck{VertexId: string("test-id")},
				&command.ComputeAck{VertexId: string("test-id"), Halted: false, AggregatedValues: make(map[string]*any.Any)},
				&command.SuperStepBarrierAck{VertexId: string("test-id")},
				&command.ComputeAck{VertexId: string("test-id"), Halted: false, AggregatedValues: make(map[string]*any.Any)},
				&command.SuperStepBarrierAck{VertexId: string("test-id")},
				&command.ComputeAck{VertexId: string("test-id"), Halted: true, AggregatedValues: make(map[string]*any.Any)},
			},
			wantComputed: 2,
			wantSentMessages: []*command.SuperStepMessage{
				{
					Uuid:         "", // ignore
					SuperStep:    0,
					SrcVertexId:  "test-id",
					DestVertexId: "dest-0",
					Message:      timestampToAny(t, 123456, 0),
				},
				{
					Uuid:         "", // ignore
					SuperStep:    1,
					SrcVertexId:  "test-id",
					DestVertexId: "dest-1",
					Message:      timestampToAny(t, 123456, 1),
				},
			},
		},
		{
			name: "fails when compute() fails",
			vertex: &MockedVertex{
				ComputeMock: func(ctx plugin.ComputeContext) error {
					return errors.New("test")
				},
				GetIDMock: func() plugin.VertexID { return "test-id" },
			},
			cmd: []proto.Message{
				&command.LoadVertex{VertexId: "test-id", PartitionId: 123},
				&command.SuperStepBarrier{},
				&command.Compute{SuperStep: 0},
			},
			wantRespond: []proto.Message{
				&command.LoadVertexAck{VertexId: "test-id", PartitionId: 123},
				&command.SuperStepBarrierAck{VertexId: "test-id"},
				nil,
			},
		},
		{
			name: "fails when receives compute before initialized",
			vertex: &MockedVertex{
				GetIDMock: func() plugin.VertexID { return "test-id" },
			},
			cmd: []proto.Message{
				&command.Compute{SuperStep: 0},
			},
			wantRespond: []proto.Message{
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			computed = 0
			context := actor.EmptyRootContext

			var childLock = &sync.Mutex{}
			var child *actor.PID
			var sentMessages []*command.SuperStepMessage
			parent := context.Spawn(actor.PropsFromFunc(func(ctx actor.Context) {
				childLock.Lock()
				defer childLock.Unlock()
				switch m := ctx.Message().(type) {
				case *actor.Started:
					logger, _ := test.NewNullLogger()
					plugin := &MockedPlugin{
						NewVertexMock: func(id plugin.VertexID) (plugin.Vertex, error) {
							return tt.vertex, nil
						},
						UnmarshalMessageMock: func(a *any.Any) (plugin.Message, error) {
							switch a.TypeUrl {
							case "com.example/test":
								return string(a.Value), nil
							case "com.github/rerorero/" + proto.MessageName(&types.Timestamp{}):
								var ts types.Timestamp
								if err := ptypes.UnmarshalAny(a, &ts); err != nil {
									t.Fatal(err)
								}
								return &ts, nil
							}
							t.Fatalf("unknown type url: %s", a.TypeUrl)
							return nil, nil
						},
						MarshalMessageMock: func(msg plugin.Message) (i *any.Any, e error) {
							switch m := msg.(type) {
							case string:
								return &any.Any{
									TypeUrl: "com.example/test",
									Value:   []byte(m),
								}, nil
							case *types.Timestamp:
								return timestampToAny(t, m.Seconds, m.Nanos), nil
							}
							t.Fatalf("unknown type: %+v", m)
							return nil, nil
						},
						GetAggregatorsMock: func() []plugin.Aggregator {
							return nil
						},
					}
					props := actor.PropsFromProducer(func() actor.Actor {
						return NewVertexActor(plugin, logger)
					})
					child = ctx.Spawn(props)
				case *command.SuperStepMessage:
					sentMessages = append(sentMessages, m)
					ctx.Respond(&command.SuperStepMessageAck{
						Uuid: m.Uuid,
					})
				case proto.Message:
					ctx.Forward(child)
				}
			}))

			for i, cmd := range tt.cmd {
				res, err := context.RequestFuture(parent, cmd, time.Second).Result()
				if (err != nil) != (tt.wantRespond[i] == nil) {
					t.Fatalf("i=%v: %v", i, err)
				}
				if diff := cmp.Diff(tt.wantRespond[i], res); diff != "" {
					t.Fatalf("unexpected respond of %d: %s", i, diff)
				}
				// send super step messages
				if msgs, ok := tt.incomingMessages[i]; ok {
					for _, msg := range msgs {
						childLock.Lock()
						ack, err := context.RequestFuture(child, msg, time.Second).Result()
						if err != nil {
							t.Fatalf("no Ack: i=%d, msg=%+v", i, *msg)
						}
						if ack.(*command.SuperStepMessageAck).Uuid != msg.Uuid {
							t.Fatalf("unexpected Ack: %+v", ack)
						}
						childLock.Unlock()
					}
				}
			}

			if computed != tt.wantComputed {
				t.Errorf("unexpected computed count: %d, %d", computed, tt.wantComputed)
			}

			if len(tt.wantSentMessages) != len(sentMessages) {
				t.Fatalf("unexpected number of messages: %d", len(sentMessages))
			}
			ignoreFields := cmpopts.IgnoreFields(command.SuperStepMessage{}, "Uuid")
			for i := range tt.wantSentMessages {
				if diff := cmp.Diff(*tt.wantSentMessages[i], *sentMessages[i], ignoreFields); diff != "" {
					t.Errorf("unexpected messages: %s", diff)
				}
			}
		})
	}
}

func timestampToAny(t *testing.T, sec int64, nano int32) *any.Any {
	ts := types.Timestamp{
		Seconds: sec,
		Nanos:   nano,
	}
	b, err := proto.Marshal(&ts)
	if err != nil {
		t.Fatal(err)
	}
	return &any.Any{
		TypeUrl: "com.github/rerorero/" + proto.MessageName(&ts),
		Value:   b,
	}
}

var aggregators = []plugin.Aggregator{
	&MockedAggregator{
		NameMock: func() string {
			return "concat"
		},
		AggregateMock: func(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) (plugin.AggregatableValue, error) {
			return v1.(string) + v2.(string), nil
		},
		MarshalValueMock: func(v plugin.AggregatableValue) (*any.Any, error) {
			return &any.Any{Value: []byte(v.(string))}, nil
		},
		UnmarshalValueMock: func(pb *any.Any) (plugin.AggregatableValue, error) {
			return string(pb.Value), nil
		},
	},
	&MockedAggregator{
		NameMock: func() string {
			return "sum"
		},
		AggregateMock: func(v1 plugin.AggregatableValue, v2 plugin.AggregatableValue) (plugin.AggregatableValue, error) {
			return v1.(uint8) + v2.(uint8), nil
		},
		MarshalValueMock: func(v plugin.AggregatableValue) (*any.Any, error) {
			return &any.Any{Value: []byte{v.(uint8)}}, nil
		},
		UnmarshalValueMock: func(pb *any.Any) (plugin.AggregatableValue, error) {
			return uint8(pb.Value[0]), nil
		},
	},
}

func Test_computeContextImpl_Aggregator(t *testing.T) {
	logger, _ := test.NewNullLogger()

	type fields struct {
		aggregatedPrevStep map[string]*any.Any
		plugin             plugin.Plugin
	}
	type args struct {
		aggregatorName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    plugin.AggregatableValue
		want1   bool
		wantErr bool
	}{
		{
			name: "get value by name",
			fields: fields{
				aggregatedPrevStep: map[string]*any.Any{
					"concat": &any.Any{Value: []byte("AA")},
				},
				plugin: &MockedPlugin{
					GetAggregatorsMock: func() []plugin.Aggregator {
						return aggregators
					},
				},
			},
			args: args{
				aggregatorName: "concat",
			},
			want:    "AA",
			want1:   true,
			wantErr: false,
		},
		{
			name: "get none value by name",
			fields: fields{
				aggregatedPrevStep: map[string]*any.Any{
					"concat": &any.Any{Value: []byte("AA")},
				},
				plugin: &MockedPlugin{
					GetAggregatorsMock: func() []plugin.Aggregator {
						return aggregators
					},
				},
			},
			args: args{
				aggregatorName: "sum",
			},
			want:    nil,
			want1:   false,
			wantErr: false,
		},
		{
			name: "unknown aggregator",
			fields: fields{
				aggregatedPrevStep: map[string]*any.Any{
					"concat": &any.Any{Value: []byte("AA")},
				},
				plugin: &MockedPlugin{
					GetAggregatorsMock: func() []plugin.Aggregator {
						return aggregators
					},
				},
			},
			args: args{
				aggregatorName: "foo",
			},
			want:    nil,
			want1:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &computeContextImpl{
				superStep: 0,
				vertexActor: &vertexActor{
					ActorUtil:             util.ActorUtil{Logger: logger},
					plugin:                tt.fields.plugin,
					aggregatedCurrentStep: nil,
				},
				aggregatedPrevStep: tt.fields.aggregatedPrevStep,
			}
			got, got1, err := c.GetAggregated(tt.args.aggregatorName)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeContextImpl.GetAggregated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeContextImpl.GetAggregated() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("computeContextImpl.GetAggregated() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_computeContextImpl_PutAggregatable(t *testing.T) {
	logger, _ := test.NewNullLogger()
	type fields struct {
		aggregatedCurrentStep map[string]*any.Any
		plugin                plugin.Plugin
	}
	type args struct {
		aggregatorName string
		v              plugin.AggregatableValue
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantValues map[string]*any.Any
	}{
		{
			name: "aggregate: append",
			fields: fields{
				aggregatedCurrentStep: make(map[string]*any.Any),
				plugin: &MockedPlugin{
					GetAggregatorsMock: func() []plugin.Aggregator {
						return aggregators
					},
				},
			},
			args: args{
				aggregatorName: "concat",
				v:              "AA",
			},
			wantErr: false,
			wantValues: map[string]*any.Any{
				"concat": &any.Any{Value: []byte("AA")},
			},
		},
		{
			name: "aggregate: reduce",
			fields: fields{
				aggregatedCurrentStep: map[string]*any.Any{
					"concat": &any.Any{Value: []byte("AA")},
					"sum":    &any.Any{Value: []byte{uint8(10)}},
				},
				plugin: &MockedPlugin{
					GetAggregatorsMock: func() []plugin.Aggregator {
						return aggregators
					},
				},
			},
			args: args{
				aggregatorName: "sum",
				v:              uint8(3),
			},
			wantErr: false,
			wantValues: map[string]*any.Any{
				"concat": &any.Any{Value: []byte("AA")},
				"sum":    &any.Any{Value: []byte{uint8(13)}},
			},
		},
		{
			name: "unknown name",
			fields: fields{
				aggregatedCurrentStep: make(map[string]*any.Any),
				plugin: &MockedPlugin{
					GetAggregatorsMock: func() []plugin.Aggregator {
						return aggregators
					},
				},
			},
			args: args{
				aggregatorName: "foooo",
				v:              "AA",
			},
			wantErr:    true,
			wantValues: make(map[string]*any.Any),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &computeContextImpl{
				vertexActor: &vertexActor{
					ActorUtil:             util.ActorUtil{Logger: logger},
					plugin:                tt.fields.plugin,
					aggregatedCurrentStep: tt.fields.aggregatedCurrentStep,
				},
			}
			if err := c.PutAggregatable(tt.args.aggregatorName, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("computeContextImpl.PutAggregatable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(c.vertexActor.aggregatedCurrentStep, tt.wantValues); diff != "" {
				t.Errorf("values different: %s", diff)
			}
		})
	}
}
