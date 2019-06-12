package worker

type ackRecorder struct {
	m map[string]struct{}
}

func (ar *ackRecorder) clear() {
	ar.m = make(map[string]struct{})
}

func (ar *ackRecorder) addToWaitList(id string) bool {
	if _, ok := ar.m[id]; ok {
		return false
	}
	ar.m[id] = struct{}{}
	return true
}

func (ar *ackRecorder) ack(id string) bool {
	if _, ok := ar.m[id]; !ok {
		return false
	}
	delete(ar.m, id)
	return true
}

func (ar *ackRecorder) hasCompleted() bool {
	return len(ar.m) == 0
}
