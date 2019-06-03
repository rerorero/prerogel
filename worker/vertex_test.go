package worker

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rerorero/prerogel/worker/command"
	"github.com/sirupsen/logrus/hooks/test"
)

func Test_vertexActor_Receive_InitVertex(t *testing.T) {
	dummyVertex := &MockedVertex{
		GetIDMock: func() VertexID {
			return VertexID("test-id")
		},
	}
	type fields struct {
		plugin Plugin
	}
	tests := []struct {
		name        string
		fields      fields
		cmd         []proto.Message
		wantRespond []proto.Message
	}{
		{
			name: "should be initialized",
			fields: fields{
				plugin: &MockedPlugin{
					NewVertexMock: func(id VertexID) Vertex {
						if id != VertexID("test-id") {
							t.Fatal("unexpected id")
						}
						return dummyVertex
					},
				},
			},
			cmd: []proto.Message{
				&command.InitVertex{
					VertexId: "test-id",
				},
			},
			wantRespond: []proto.Message{
				&command.InitVertexAck{
					VertexId: "test-id",
				},
			},
		},
		{
			name: "fails if already initialized",
			fields: fields{
				plugin: &MockedPlugin{
					NewVertexMock: func(id VertexID) Vertex {
						return dummyVertex
					},
				},
			},
			cmd: []proto.Message{
				&command.InitVertex{
					VertexId: "test-id",
				},
				&command.InitVertex{},
			},
			wantRespond: []proto.Message{
				&command.InitVertexAck{
					VertexId: "test-id",
				},
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := test.NewNullLogger()

			context := actor.EmptyRootContext
			props := actor.PropsFromProducer(func() actor.Actor {
				return NewVertexActor(tt.fields.plugin, logger)
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

func Test_vertexActor_Receive_LoadVertex(t *testing.T) {
	var loaded int
	tests := []struct {
		name        string
		vertex      Vertex
		cmd         []proto.Message
		wantRespond []proto.Message
		wantLoaded  int
	}{
		{
			name: "load",
			vertex: &MockedVertex{
				LoadMock: func() error {
					loaded++
					return nil
				},
				GetIDMock: func() VertexID {
					return "test-id"
				},
			},
			cmd: []proto.Message{
				&command.InitVertex{VertexId: "test-id"},
				&command.LoadVertex{},
			},
			wantRespond: []proto.Message{
				&command.InitVertexAck{
					VertexId: "test-id",
				},
				&command.LoadVertexAck{
					VertexId: "test-id",
				},
			},
			wantLoaded: 1,
		},
		{
			name: "fail if has not initialized yet",
			vertex: &MockedVertex{
				LoadMock: func() error {
					loaded++
					return nil
				},
				GetIDMock: func() VertexID {
					return "test-id"
				},
			},
			cmd: []proto.Message{
				&command.LoadVertex{},
			},
			wantRespond: []proto.Message{
				nil,
			},
			wantLoaded: 0,
		},
		{
			name: "fails when Load() returns error",
			vertex: &MockedVertex{
				LoadMock: func() error {
					loaded++
					return errors.New("test")
				},
				GetIDMock: func() VertexID {
					return "test-id"
				},
			},
			cmd: []proto.Message{
				&command.InitVertex{VertexId: "test-id"},
				&command.LoadVertex{},
			},
			wantRespond: []proto.Message{
				&command.InitVertexAck{
					VertexId: "test-id",
				},
				nil,
			},
			wantLoaded: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loaded = 0
			logger, _ := test.NewNullLogger()
			plugin := &MockedPlugin{
				NewVertexMock: func(id VertexID) Vertex {
					return tt.vertex
				},
			}

			context := actor.EmptyRootContext
			props := actor.PropsFromProducer(func() actor.Actor {
				return NewVertexActor(plugin, logger)
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
			if loaded != tt.wantLoaded {
				t.Errorf("unexpected loaded count: %d, %d", loaded, tt.wantLoaded)
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
		Message:      &types.Any{Value: []byte("test1")},
	}
	msg2 := &command.SuperStepMessage{
		Uuid:         "b",
		SuperStep:    0,
		SrcVertexId:  "bar",
		DestVertexId: "test-id",
		Message:      &types.Any{Value: []byte("test2")},
	}

	tests := []struct {
		name             string
		vertex           Vertex
		cmd              []proto.Message
		incomingMessages map[int][]*command.SuperStepMessage
		wantRespond      []proto.Message
		wantComputed     int
		wantSentMessages []*command.SuperStepMessage
	}{
		{
			name: "compute 0 then halt",
			vertex: &MockedVertex{
				ComputeMock: func(ctx ComputeContext) error {
					if ctx.SuperStep() != 0 {
						t.Fatal("unexpected step")
					}
					if ctx.ReceivedMessages() != nil {
						t.Fatal("unexpected received messages")
					}
					computed++
					return nil
				},
				GetIDMock: func() VertexID {
					return "test-id"
				},
			},
			cmd: []proto.Message{
				&command.InitVertex{VertexId: "test-id"},
				&command.Compute{SuperStep: 0},
				&command.Compute{SuperStep: 1},
			},
			wantRespond: []proto.Message{
				&command.InitVertexAck{
					VertexId: "test-id",
				},
				&command.ComputeAck{
					VertexId: string("test-id"),
					Halted:   false,
				},
				&command.ComputeAck{
					VertexId: string("test-id"),
					Halted:   true,
				},
			},
			wantComputed:     1,
			wantSentMessages: nil,
		},
		{
			name: "receive messages and continue to compute",
			vertex: &MockedVertex{
				ComputeMock: func(ctx ComputeContext) error {
					if ctx.SuperStep() != uint64(computed) {
						t.Fatalf("unexpected step: %v, %v", ctx.SuperStep(), computed)
					}
					if ctx.SuperStep() == 1 {
						if diff := cmp.Diff([]Message{msg1.Message, msg2.Message}, ctx.ReceivedMessages()); diff != "" {
							t.Fatalf("unexpected received messages: %s", diff)
						}
					} else if ctx.ReceivedMessages() != nil {
						t.Fatal("unexpected received messages")
					}
					computed++
					return nil
				},
				GetIDMock: func() VertexID { return "test-id" },
			},
			cmd: []proto.Message{
				&command.InitVertex{VertexId: "test-id"},
				&command.Compute{SuperStep: 0},
				&command.Compute{SuperStep: 1},
				&command.Compute{SuperStep: 2},
			},
			incomingMessages: map[int][]*command.SuperStepMessage{
				1: {msg1, msg2}, // sent between Compute(step0) and Compute(step1)
			},
			wantRespond: []proto.Message{
				&command.InitVertexAck{
					VertexId: "test-id",
				},
				&command.ComputeAck{
					VertexId: string("test-id"),
					Halted:   false,
				},
				&command.ComputeAck{
					VertexId: string("test-id"),
					Halted:   false,
				},
				&command.ComputeAck{
					VertexId: string("test-id"),
					Halted:   true,
				},
			},
			wantComputed:     2,
			wantSentMessages: nil,
		},
		{
			name: "send messages to other vertices",
			vertex: &MockedVertex{
				ComputeMock: func(ctx ComputeContext) error {
					msg := &types.Timestamp{
						Seconds: 123456,
						Nanos:   int32(ctx.SuperStep()),
					}
					if err := ctx.SendMessageTo(VertexID(fmt.Sprintf("dest-%v", ctx.SuperStep())), msg); err != nil {
						t.Fatal(err)
					}
					computed++
					return nil
				},
				GetIDMock: func() VertexID { return "test-id" },
			},
			cmd: []proto.Message{
				&command.InitVertex{VertexId: "test-id"},
				&command.Compute{SuperStep: 0},
				&command.Compute{SuperStep: 1},
				&command.Compute{SuperStep: 2},
			},
			incomingMessages: map[int][]*command.SuperStepMessage{
				1: {msg1}, // sent between Compute(step0) and Compute(step1)
			},
			wantRespond: []proto.Message{
				&command.InitVertexAck{
					VertexId: "test-id",
				},
				&command.ComputeAck{
					VertexId: string("test-id"),
					Halted:   false,
				},
				&command.ComputeAck{
					VertexId: string("test-id"),
					Halted:   false,
				},
				&command.ComputeAck{
					VertexId: string("test-id"),
					Halted:   true,
				},
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
				ComputeMock: func(ctx ComputeContext) error {
					return errors.New("test")
				},
				GetIDMock: func() VertexID { return "test-id" },
			},
			cmd: []proto.Message{
				&command.InitVertex{VertexId: "test-id"},
				&command.Compute{SuperStep: 0},
			},
			wantRespond: []proto.Message{
				&command.InitVertexAck{
					VertexId: "test-id",
				},
				nil,
			},
		},
		{
			name: "fails when receives compute before initialized",
			vertex: &MockedVertex{
				GetIDMock: func() VertexID { return "test-id" },
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
						NewVertexMock: func(id VertexID) Vertex {
							return tt.vertex
						},
					}
					props := actor.PropsFromProducer(func() actor.Actor {
						return NewVertexActor(plugin, logger)
					})
					child = ctx.Spawn(props)
				case *command.SuperStepMessage:
					sentMessages = append(sentMessages, m)
				case proto.Message:
					ctx.Forward(child)
				}
			}))

			for i, cmd := range tt.cmd {
				res, err := context.RequestFuture(parent, cmd, time.Second).Result()
				if (err != nil) != (tt.wantRespond[i] == nil) {
					t.Fatal(err)
				}
				if diff := cmp.Diff(tt.wantRespond[i], res); diff != "" {
					t.Fatalf("unexpected respond of %d: %s", i, diff)
				}
				// send super step messages
				msgs, ok := tt.incomingMessages[i]
				if ok {
					for _, msg := range msgs {
						childLock.Lock()
						ack, err := context.RequestFuture(child, msg, time.Second).Result()
						if err != nil {
							t.Fatalf("no ack: i=%d, msg=%+v", i, *msg)
						}
						if ack.(*command.SuperStepMessageAck).Uuid != msg.Uuid {
							t.Fatalf("unexpected ack: %+v", ack)
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

func timestampToAny(t *testing.T, sec int64, nano int32) *types.Any {
	a, err := types.MarshalAny(&types.Timestamp{
		Seconds: sec,
		Nanos:   nano,
	})
	if err != nil {
		t.Fatal(err)
	}
	return a
}
