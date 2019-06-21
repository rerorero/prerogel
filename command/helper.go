package command

// StatsCompleted returns if CoordinatorStatsAck shows processing has finished or not
func (s *CoordinatorStatsAck) StatsCompleted() bool {
	// Why I check the number of messages sent is that the number of actives is often incorrect.
	// Vertex actor returns its active state with ComputeAck, but then it may receives a message until the next superstep is started
	return s.SuperStep > 0 && s.NrOfActiveVertex == 0 && s.NrOfSentMessages == 0
}
