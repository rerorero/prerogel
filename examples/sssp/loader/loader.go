package loader

// Vert is vertex
type Vert struct {
	ID        string
	Outgoings map[string]uint32
}

// Loader loads graph
type Loader interface {
	LoadPartition(partitionID uint64, numOfPartitions uint64) ([]*Vert, error)
}
