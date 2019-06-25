package worker

import "github.com/rerorero/prerogel/aggregator"

const (
	// WorkerActorKind is actor kind of worker
	WorkerActorKind = "worker"

	// CoordinatorActorID is actor id of coordinator
	CoordinatorActorID = "coordinator"

	// VertexStatsName is aggregator name of VertexStatsAggregator
	VertexStatsName = "prerogel/vertex-stats"
)

// VertexStatsAggregatorInstance is singleton
var vertexStatsAggregatorInstance = &aggregator.VertexStatsAggregator{
	AggName: VertexStatsName,
}
