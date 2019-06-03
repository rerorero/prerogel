package worker

import (
	"errors"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"

	"github.com/gogo/protobuf/proto"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/google/go-cmp/cmp"
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
				&command.InitVertexCompleted{
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
				&command.InitVertexCompleted{
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
				&command.InitVertexCompleted{
					VertexId: "test-id",
				},
				&command.LoadVertexCompleted{
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
				&command.InitVertexCompleted{
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
				&command.InitVertexCompleted{
					VertexId: "test-id",
				},
				&command.ComputeCompleted{
					VertexId: string("test-id"),
					Halted:   false,
				},
				&command.ComputeCompleted{
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
					if ctx.SuperStep() != 0 {
						t.Fatal("unexpected step")
					}
					if ctx.ReceivedMessages() != nil {
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
				0: {
					{
						SuperStep:    0,
						DestVertexId: "test-id",
						Message:      &types.Any{Value: []byte("test1")},
					},
					{
						SuperStep:    0,
						DestVertexId: "test-id",
						Message:      &types.Any{Value: []byte("test2")},
					},
				},
			},
			wantRespond: []proto.Message{
				&command.InitVertexCompleted{
					VertexId: "test-id",
				},
				&command.ComputeCompleted{
					VertexId: string("test-id"),
					Halted:   false,
				},
				&command.ComputeCompleted{
					VertexId: string("test-id"),
					Halted:   true,
				},
			},
			wantComputed:     1,
			wantSentMessages: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			var child *actor.PID
			var sentMessages []*command.SuperStepMessage
			parent := context.Spawn(actor.PropsFromFunc(func(ctx actor.Context) {
				switch m := ctx.Message().(type) {
				case *actor.Started:
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
					t.Errorf("unexpected respond: %s", diff)
				}
				// send super step messages
				msgs := tt.incomingMessages[i]
				for _, msg := range msgs {
					context.Send(parent, msg)
				}
			}
			if computed != tt.wantComputed {
				t.Errorf("unexpected computed count: %d, %d", computed, tt.wantComputed)
			}
			if diff := cmp.Diff(tt.wantSentMessages, sentMessages); diff != "" {
				t.Errorf("unexpected messages: %s", diff)
			}
		})
	}
}
