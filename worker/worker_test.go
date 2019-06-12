package worker

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rerorero/prerogel/worker/command"
)

func Test_superStepMsgBuf_add_remove(t *testing.T) {
	buf := newSuperStepMsgBuf(nil)

	// add
	m1 := &command.SuperStepMessage{
		Uuid:           "uuid1",
		SrcVertexId:    "s1",
		SrcPartitionId: 1,
		DestVertexId:   "d1",
		Message:        anyOf("m1"),
	}
	m2 := &command.SuperStepMessage{
		Uuid:           "uuid2",
		SrcVertexId:    "s2",
		SrcPartitionId: 1,
		DestVertexId:   "d1",
		Message:        anyOf("m2"),
	}
	m3 := &command.SuperStepMessage{
		Uuid:           "uuid3",
		SrcVertexId:    "s3",
		SrcPartitionId: 1,
		DestVertexId:   "d2",
		Message:        anyOf("m3"),
	}

	// add
	buf.add(m1)
	buf.add(m2)
	buf.add(m3)
	expected := map[VertexID][]*command.SuperStepMessage{
		VertexID("d1"): {m1, m2},
		VertexID("d2"): {m3},
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
	expected = map[VertexID][]*command.SuperStepMessage{
		VertexID("d1"): {m1},
		VertexID("d2"): {m3},
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
	expected = map[VertexID][]*command.SuperStepMessage{
		VertexID("d2"): {m3},
	}
	if diff := cmp.Diff(expected, buf.buf); diff != "" {
		t.Fatalf("not match: %s", diff)
	}
	if buf.numOfMessage() != 1 {
		t.Fatal("unexpected number")
	}

	// Clear
	buf.clear()
	expected = map[VertexID][]*command.SuperStepMessage{}
	if diff := cmp.Diff(expected, buf.buf); diff != "" {
		t.Fatalf("not match: %s", diff)
	}
	if buf.numOfMessage() != 0 {
		t.Fatal("unexpected number")
	}
}

func Test_superStepMsgBuf_combine(t *testing.T) {
	buf := newSuperStepMsgBuf(&MockedPlugin{
		MarshalMessageMock: func(msg Message) (*any.Any, error) {
			return anyOf(msg.(string)), nil
		},
		UnmarshalMessageMock: func(a *any.Any) (Message, error) {
			return string(a.Value), nil
		},
		GetCombinerMock: func() func(VertexID, []Message) ([]Message, error) {
			return func(id VertexID, msgs []Message) ([]Message, error) {
				// choose longest string
				longest := msgs[0].(string)
				for _, m := range msgs {
					if len(m.(string)) > len(longest) {
						longest = m.(string)
					}
				}
				return []Message{longest}, nil
			}
		},
	})

	m1 := &command.SuperStepMessage{
		Uuid:           "uuid1",
		SrcVertexId:    "s1",
		SrcPartitionId: 1,
		DestVertexId:   "d1",
		Message:        anyOf("m1"),
	}
	m2 := &command.SuperStepMessage{
		Uuid:           "uuid2",
		SrcVertexId:    "s2",
		SrcPartitionId: 1,
		DestVertexId:   "d1",
		Message:        anyOf("m2_middle"),
	}
	m3 := &command.SuperStepMessage{
		Uuid:           "uuid3",
		SrcVertexId:    "s2",
		SrcPartitionId: 1,
		DestVertexId:   "d1",
		Message:        anyOf("m2_looooooooooooooooong"),
	}
	m4 := &command.SuperStepMessage{
		Uuid:           "uuid4",
		SrcVertexId:    "s3",
		SrcPartitionId: 1,
		DestVertexId:   "d2",
		Message:        anyOf("m4_long"),
	}
	m5 := &command.SuperStepMessage{
		Uuid:           "uuid5",
		SrcVertexId:    "s3",
		SrcPartitionId: 1,
		DestVertexId:   "d2",
		Message:        anyOf("m5"),
	}
	m6 := &command.SuperStepMessage{
		Uuid:           "uuid6",
		SrcVertexId:    "s1",
		SrcPartitionId: 1,
		DestVertexId:   "d3",
		Message:        anyOf("m6"),
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

	expected := map[VertexID][]*command.SuperStepMessage{
		VertexID("d1"): {&command.SuperStepMessage{
			Uuid:           "",
			SrcVertexId:    "",
			SrcPartitionId: 1,
			DestVertexId:   "d1",
			Message:        anyOf("m2_looooooooooooooooong"),
		}},
		VertexID("d2"): {&command.SuperStepMessage{
			Uuid:           "",
			SrcVertexId:    "",
			SrcPartitionId: 1,
			DestVertexId:   "d2",
			Message:        anyOf("m4_long"),
		}},
		VertexID("d3"): {m6},
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
