package command

// StatsCompleted returns if CoordinatorStatsAck shows processing has finished or not
func (s *CoordinatorStatsAck) StatsCompleted() bool {
	// Why I check the number of messages sent is that the number of actives is often incorrect.
	// Vertex actor returns its active state with ComputeAck, but then it may receives a message until the next superstep is started
	return s.SuperStep > 0 && s.NrOfActiveVertex == 0 && s.NrOfSentMessages == 0
}

// FindWoerkerInfoByPartition returns worker info that owns partition
func (info *ClusterInfo) FindWoerkerInfoByPartition(partitionID uint64) *ClusterInfo_WorkerInfo {
	for _, info := range info.WorkerInfo {
		for _, p := range info.Partitions {
			if p == partitionID {
				return info
			}
		}
	}
	return nil
}

// NumOfPartitions returns number of partitions
func (info *ClusterInfo) NumOfPartitions() uint64 {
	var size uint64
	for _, i := range info.WorkerInfo {
		size += uint64(len(i.Partitions))
	}
	return size
}
