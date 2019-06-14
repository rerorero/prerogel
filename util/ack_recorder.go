package util

// AckRecorder helps to manage state of waiting Ack message
type AckRecorder struct {
	m map[string]struct{}
}

// Clear clears waiting list
func (ar *AckRecorder) Clear() {
	ar.m = make(map[string]struct{})
}

// AddToWaitList adds id into waiting list
func (ar *AckRecorder) AddToWaitList(id string) bool {
	if _, ok := ar.m[id]; ok {
		return false
	}
	ar.m[id] = struct{}{}
	return true
}

// Ack removes id from waiting list
func (ar *AckRecorder) Ack(id string) bool {
	if _, ok := ar.m[id]; !ok {
		return false
	}
	delete(ar.m, id)
	return true
}

// HasCompleted checks if waiting list is empty
func (ar *AckRecorder) HasCompleted() bool {
	return len(ar.m) == 0
}

// Size returns size
func (ar *AckRecorder) Size() int {
	return len(ar.m)
}
