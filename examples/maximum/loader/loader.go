package loader

// Loader loads graph
type Loader interface {
	LoadVertex(id string) (uint32, []string, error)
	LoadPartition(partitionID uint64, numOfPartitions uint64) (map[string]uint32, map[string][]string, error)
}
