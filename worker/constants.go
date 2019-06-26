package worker

import (
	"github.com/rerorero/prerogel/aggregator"
	"github.com/rerorero/prerogel/plugin"
)

const (
	// WorkerActorKind is actor kind of worker
	WorkerActorKind = "worker"

	// CoordinatorActorID is actor id of coordinator
	CoordinatorActorID = "coordinator"

	// VertexStatsName is aggregator name of VertexStatsAggregator
	VertexStatsName = "prerogel/vertex-stats"
)

var (
	// vertexStatsAggregatorInstance is singleton
	vertexStatsAggregatorInstance = &aggregator.VertexStatsAggregator{
		AggName: VertexStatsName,
	}

	systemAggregator = []plugin.Aggregator{
		vertexStatsAggregatorInstance,
	}
)
