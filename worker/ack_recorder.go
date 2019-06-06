package worker

func (ar *ackRecorder) clear() {
	ar.m = make(map[VertexID]struct{})
}

func (ar *ackRecorder) setExpected(map[VertexID]) {
	ar.m = make(map[VertexID]struct{})
}

// ack marks that id has received ack, and return if already received ack or not as bool
func (ar *ackRecorder) ack(id VertexID) bool {
	if _, ok := ar.m[id]; ok {
		return true
	}
	ar.m[id] = struct{}{}
	return false
}

func (ar *ackRecorder) hasCompleted(expected int) bool {
	return expected == len(ar.m)
}
