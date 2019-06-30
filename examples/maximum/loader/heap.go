package loader

import (
	"fmt"

	"github.com/rerorero/prerogel/plugin"
)

// HeapLoader load vertex from heap
type HeapLoader struct {
	Vertices map[string]uint32
	Edges    [][]string
}

// LoadVertex reads vertex and outgoing edges
func (l *HeapLoader) LoadVertex(id string) (uint32, []string, error) {
	val, ok := l.Vertices[id]
	if !ok {
		return 0, nil, fmt.Errorf("id %s doesn't exist", id)
	}

	return val, l.edgeOf(id), nil
}

// LoadPartition loads vertices of partition
func (l *HeapLoader) LoadPartition(partitionID uint64, numOfPartitions uint64) (map[string]uint32, map[string][]string, error) {
	vertex := make(map[string]uint32)
	outgoings := make(map[string][]string)

	for id, v := range l.Vertices {
		// this is so bad implementation
		partition, err := plugin.HashPartition(plugin.VertexID(id), numOfPartitions)
		if err != nil {
			return nil, nil, err
		}
		if partition == partitionID {
			vertex[id] = v
			outgoings[id] = l.edgeOf(id)
		}
	}

	return vertex, outgoings, nil
}

func (l *HeapLoader) edgeOf(id string) []string {
	var edges []string
	for _, e := range l.Edges {
		if e[0] == id {
			edges = append(edges, e[1])
		}
	}
	return edges
}
