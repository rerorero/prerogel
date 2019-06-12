package util

import "testing"

func Test_ackRecorder_clear(t *testing.T) {
	ar := &AckRecorder{}
	ar.Clear()
	if len(ar.m) != 0 {
		t.Fatal("not empty")
	}
	if !ar.HasCompleted() {
		t.Fatal("completed")
	}

	if !ar.AddToWaitList("test") {
		t.Fatal("dup")
	}
	if ar.HasCompleted() {
		t.Fatal("completed")
	}
	if ar.AddToWaitList("test") {
		t.Fatal("dup")
	}

	if !ar.AddToWaitList("test2") {
		t.Fatal("dup")
	}
	if !ar.Ack("test") {
		t.Fatal("dup")
	}
	if ar.Ack("test") {
		t.Fatal("dup")
	}
	if ar.HasCompleted() {
		t.Fatal("completed")
	}

	if !ar.Ack("test2") {
		t.Fatal("dup")
	}
	if !ar.HasCompleted() {
		t.Fatal("completed")
	}
}
