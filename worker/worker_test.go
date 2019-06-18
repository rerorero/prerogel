package worker

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/plugin"
	"github.com/rerorero/prerogel/util"
	"github.com/sirupsen/logrus/hooks/test"
)

func Test_superStepMsgBuf_add_remove(t *testing.T) {
	buf := newSuperStepMsgBuf(nil)

	// add
	m1 := &command.SuperStepMessage{
		Uuid:         "uuid1",
		SrcVertexId:  "s1",
		DestVertexId: "d1",
		Message:      anyOf("m1"),
	}
	m2 := &command.SuperStepMessage{
		Uuid:         "uuid2",
		SrcVertexId:  "s2",
		DestVertexId: "d1",
		Message:      anyOf("m2"),
	}
	m3 := &command.SuperStepMessage{
		Uuid:         "uuid3",
		SrcVertexId:  "s3",
		DestVertexId: "d2",
		Message:      anyOf("m3"),
	}

	// add
	buf.add(m1)
	buf.add(m2)
	buf.add(m3)
	expected := map[plugin.VertexID][]*command.SuperStepMessage{
		plugin.VertexID("d1"): {m1, m2},
		plugin.VertexID("d2"): {m3},
	}
	if diff := cmp.Diff(expected, buf.buf); diff != "" {
		t.Fatalf("not match: %s", diff)
	}
	if buf.numOfMessage() != 3 {
		t.Fatal("unexpected number")
	}

	// remove 1
	buf.remove(&command.SuperStepMessageAck{
		Uuid: "uuid2",
	})
	expected = map[plugin.VertexID][]*command.SuperStepMessage{
		plugin.VertexID("d1"): {m1},
		plugin.VertexID("d2"): {m3},
	}
	if diff := cmp.Diff(expected, buf.buf); diff != "" {
		t.Fatalf("not match: %s", diff)
	}
	if buf.numOfMessage() != 2 {
		t.Fatal("unexpected number")
	}

	// remove 2
	buf.remove(&command.SuperStepMessageAck{
		Uuid: "uuid1",
	})
	expected = map[plugin.VertexID][]*command.SuperStepMessage{
		plugin.VertexID("d2"): {m3},
	}
	if diff := cmp.Diff(expected, buf.buf); diff != "" {
		t.Fatalf("not match: %s", diff)
	}
	if buf.numOfMessage() != 1 {
		t.Fatal("unexpected number")
	}

	// Clear
	buf.clear()
	expected = map[plugin.VertexID][]*command.SuperStepMessage{}
	if diff := cmp.Diff(expected, buf.buf); diff != "" {
		t.Fatalf("not match: %s", diff)
	}
	if buf.numOfMessage() != 0 {
		t.Fatal("unexpected number")
	}
}

func Test_superStepMsgBuf_combine(t *testing.T) {
	buf := newSuperStepMsgBuf(&MockedPlugin{
		MarshalMessageMock: func(msg plugin.Message) (*any.Any, error) {
			return anyOf(msg.(string)), nil
		},
		UnmarshalMessageMock: func(a *any.Any) (plugin.Message, error) {
			return string(a.Value), nil
		},
		GetCombinerMock: func() func(plugin.VertexID, []plugin.Message) ([]plugin.Message, error) {
			return func(id plugin.VertexID, msgs []plugin.Message) ([]plugin.Message, error) {
				// choose longest string
				longest := msgs[0].(string)
				for _, m := range msgs {
					if len(m.(string)) > len(longest) {
						longest = m.(string)
					}
				}
				return []plugin.Message{longest}, nil
			}
		},
	})

	m1 := &command.SuperStepMessage{
		Uuid:         "uuid1",
		SrcVertexId:  "s1",
		DestVertexId: "d1",
		Message:      anyOf("m1"),
	}
	m2 := &command.SuperStepMessage{
		Uuid:         "uuid2",
		SrcVertexId:  "s2",
		DestVertexId: "d1",
		Message:      anyOf("m2_middle"),
	}
	m3 := &command.SuperStepMessage{
		Uuid:         "uuid3",
		SrcVertexId:  "s2",
		DestVertexId: "d1",
		Message:      anyOf("m2_looooooooooooooooong"),
	}
	m4 := &command.SuperStepMessage{
		Uuid:         "uuid4",
		SrcVertexId:  "s3",
		DestVertexId: "d2",
		Message:      anyOf("m4_long"),
	}
	m5 := &command.SuperStepMessage{
		Uuid:         "uuid5",
		SrcVertexId:  "s3",
		DestVertexId: "d2",
		Message:      anyOf("m5"),
	}
	m6 := &command.SuperStepMessage{
		Uuid:         "uuid6",
		SrcVertexId:  "s1",
		DestVertexId: "d3",
		Message:      anyOf("m6"),
	}

	buf.add(m1)
	buf.add(m2)
	buf.add(m3)
	buf.add(m4)
	buf.add(m5)
	buf.add(m6)
	if err := buf.combine(); err != nil {
		t.Fatal(err)
	}

	expected := map[plugin.VertexID][]*command.SuperStepMessage{
		plugin.VertexID("d1"): {&command.SuperStepMessage{
			Uuid:         "",
			SrcVertexId:  "",
			DestVertexId: "d1",
			Message:      anyOf("m2_looooooooooooooooong"),
		}},
		plugin.VertexID("d2"): {&command.SuperStepMessage{
			Uuid:         "",
			SrcVertexId:  "",
			DestVertexId: "d2",
			Message:      anyOf("m4_long"),
		}},
		plugin.VertexID("d3"): {m6},
	}

	ignoreFields := cmpopts.IgnoreFields(command.SuperStepMessage{}, "Uuid")
	if diff := cmp.Diff(expected, buf.buf, ignoreFields); diff != "" {
		t.Fatalf("not match: %s", diff)
	}
	if buf.numOfMessage() != 3 {
		t.Fatal("unexpected number")
	}
}

func anyOf(s string) *any.Any {
	return &any.Any{
		Value: []byte(s),
	}
}

func TestNewWorkerActor_routesMessages(t *testing.T) {
	var called int32
	messageAck := make(map[int]int)
	messageAckMux := &sync.Mutex{}
	logger, _ := test.NewNullLogger()

	plugin := &MockedPlugin{
		PartitionMock: func(vid plugin.VertexID, numOfPartitions uint64) (uint64, error) {
			id := string(vid)
			if numOfPartitions != 6 {
				return 0, errors.New("plugin failure: unexpected numOfPartitions")
			}
			i, err := strconv.Atoi(id[len(id)-1:])
			if err != nil {
				t.Fatal(err)
				return 0, errors.Wrap(err, "plugin failure: unexpected numOfPartitions")
			}
			return uint64(i), nil
		},
		MarshalMessageMock: func(msg plugin.Message) (*any.Any, error) {
			return msg.(*any.Any), nil
		},
		UnmarshalMessageMock: func(a *any.Any) (plugin.Message, error) {
			return a, nil
		},
		GetCombinerMock: func() func(plugin.VertexID, []plugin.Message) ([]plugin.Message, error) {
			return func(id plugin.VertexID, messages []plugin.Message) ([]plugin.Message, error) {
				return messages, nil
			}
		},
	}

	partitions := []uint64{1, 2, 3}

	partitionProps := actor.PropsFromFunc(func(c actor.Context) {
		switch cmd := c.Message().(type) {
		case *command.InitPartition:
			c.Send(c.Parent(), &command.InitPartitionAck{PartitionId: cmd.PartitionId})
		case *command.SuperStepBarrier:
			i := atomic.AddInt32(&called, 1)
			c.Send(c.Parent(), &command.SuperStepBarrierPartitionAck{PartitionId: partitions[i-1]})
		case *command.Compute:
			i := atomic.AddInt32(&called, 1)
			// internal message
			c.Request(c.Parent(), &command.SuperStepMessage{
				Uuid:         fmt.Sprintf("uuid-internal-%v", partitions[i-1]),
				SuperStep:    cmd.SuperStep,
				SrcVertexId:  "vertex-dummy-1",
				DestVertexId: fmt.Sprintf("dest-internal-%v", partitions[i-1]),
				Message:      anyOf("message-from-me"),
			})
			// external message
			c.Request(c.Parent(), &command.SuperStepMessage{
				Uuid:         fmt.Sprintf("uuid-external-%v", partitions[i-1]),
				SuperStep:    cmd.SuperStep,
				SrcVertexId:  "vertex-dummy-2",
				DestVertexId: fmt.Sprintf("dest-external-%v", partitions[i-1]+3),
				Message:      anyOf("message-from-me"),
			})
		case *command.SuperStepMessage:
			c.Respond(&command.SuperStepMessageAck{
				Uuid: cmd.Uuid,
			})
		case *command.SuperStepMessageAck:
			id, err := strconv.Atoi(cmd.Uuid[len(cmd.Uuid)-1:])
			if err != nil {
				t.Fatal(err)
			}
			messageAckMux.Lock()
			defer messageAckMux.Unlock()
			switch {
			case strings.Contains(cmd.Uuid, "uuid-internal-"):
				messageAck[id] = messageAck[id] + 1
			case strings.Contains(cmd.Uuid, "uuid-external-"):
				messageAck[id] = messageAck[id] + 1
			default:
				t.Fatalf("unexpected ack: %#v", cmd)
			}
			if messageAck[id] == 2 {
				c.Send(c.Parent(), &command.ComputePartitionAck{PartitionId: partitions[id-1]})
			}
		}
	})

	workerProps := actor.PropsFromProducer(func() actor.Actor {
		return NewWorkerActor(plugin, partitionProps, logger)
	})
	context := actor.EmptyRootContext
	computeAckCh := make(chan *command.ComputeWorkerAck, 1)
	proxy := util.NewActorProxy(context, workerProps, func(ctx actor.Context) {
		switch cmd := ctx.Message().(type) {
		case *command.ComputeWorkerAck:
			computeAckCh <- cmd
		}
	})

	extMessageAckCh := make(chan *command.SuperStepMessageAck, 1)
	defer close(extMessageAckCh)
	otherWorkerMock := actor.PropsFromFunc(func(c actor.Context) {
		switch cmd := c.Message().(type) {
		case *command.SuperStepMessage:
			c.Respond(&command.SuperStepMessageAck{
				Uuid: cmd.Uuid,
			})
		case *command.SuperStepMessageAck:
			extMessageAckCh <- cmd
		case string:
			// external message
			c.Request(proxy.Underlying(), &command.SuperStepMessage{
				Uuid:         "uuid-ext-worker",
				SuperStep:    1,
				SrcVertexId:  "src-" + cmd,
				DestVertexId: "dest-internal-2",
				Message:      anyOf(cmd),
			})
		}
	})
	other1 := context.Spawn(otherWorkerMock)

	// move state forward
	called = 0
	if _, err := proxy.SendAndAwait(context, &command.InitWorker{
		Partitions: partitions,
	}, &command.InitWorkerAck{}, time.Second); err != nil {
		t.Fatal(err)
	}

	ci := &command.ClusterInfo{
		WorkerInfo: []*command.ClusterInfo_WorkerInfo{
			{
				WorkerPid:  proxy.Underlying(),
				Partitions: partitions,
			},
			{
				WorkerPid:  other1,
				Partitions: []uint64{4, 5, 6},
			},
		},
	}

	// step 0
	called = 0
	if _, err := proxy.SendAndAwait(context, &command.SuperStepBarrier{ClusterInfo: ci}, &command.SuperStepBarrierWorkerAck{}, time.Second); err != nil {
		t.Fatal(err)
	}
	called = 0
	messageAck = make(map[int]int)
	proxy.Send(context, &command.Compute{SuperStep: 0})
	<-computeAckCh

	// step 1
	called = 0
	if _, err := proxy.SendAndAwait(context, &command.SuperStepBarrier{ClusterInfo: ci}, &command.SuperStepBarrierWorkerAck{}, time.Second); err != nil {
		t.Fatal(err)
	}
	called = 0
	messageAck = make(map[int]int)
	proxy.Send(context, &command.Compute{SuperStep: 1})
	context.Send(other1, "ext-5")
	<-computeAckCh
	ack := <-extMessageAckCh
	if ack.Uuid != "uuid-ext-worker" {
		t.Fatal("unexpected ack")
	}

	// step 2
	called = 0
	if _, err := proxy.SendAndAwait(context, &command.SuperStepBarrier{ClusterInfo: ci}, &command.SuperStepBarrierWorkerAck{}, time.Second); err != nil {
		t.Fatal(err)
	}
	called = 0
	messageAck = make(map[int]int)
	proxy.Send(context, &command.Compute{SuperStep: 2})
	<-computeAckCh
}
