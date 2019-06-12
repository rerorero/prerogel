package worker

import "testing"

func Test_ackRecorder_clear(t *testing.T) {
	ar := &ackRecorder{}
	ar.clear()
	if len(ar.m) != 0 {
		t.Fatal("not empty")
	}
	if !ar.hasCompleted() {
		t.Fatal("completed")
	}

	if !ar.addToWaitList("test") {
		t.Fatal("dup")
	}
	if ar.hasCompleted() {
		t.Fatal("completed")
	}
	if ar.addToWaitList("test") {
		t.Fatal("dup")
	}

	if !ar.addToWaitList("test2") {
		t.Fatal("dup")
	}
	if !ar.ack("test") {
		t.Fatal("dup")
	}
	if ar.ack("test") {
		t.Fatal("dup")
	}
	if ar.hasCompleted() {
		t.Fatal("completed")
	}

	if !ar.ack("test2") {
		t.Fatal("dup")
	}
	if !ar.hasCompleted() {
		t.Fatal("completed")
	}
}
