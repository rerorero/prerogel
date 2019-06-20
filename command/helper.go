package command

// StatsCompleted returns if CoordinatorStatsAck shows processing has finished or not
func StatsCompleted(s *CoordinatorStatsAck) bool {
	return s.SuperStep > 0 && s.NrOfActiveVertex == 0
}
