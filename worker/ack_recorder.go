package worker

type ackRecorder struct {
	m map[string]struct{}
}

func (ar *ackRecorder) clear() {
	ar.m = make(map[string]struct{})
}

func (ar *ackRecorder) addToWaitList(id string) bool {
	if _, ok := ar.m[id]; ok {
		return true
	}
	ar.m[id] = struct{}{}
	return false
}

func (ar *ackRecorder) ack(id string) bool {
	if _, ok := ar.m[id]; !ok {
		return true
	}
	delete(ar.m, id)
	return false
}

func (ar *ackRecorder) hasCompleted() bool {
	return len(ar.m) == 0
}
