package loader

import "fmt"

// HeapLoader load vertex from heap
type HeapLoader struct {
	Vertices map[string]uint32
	Edges    [][]string
}

// Load reads vertex and outgoing edges
func (l *HeapLoader) Load(id string) (uint32, []string, error) {
	val, ok := l.Vertices[id]
	if !ok {
		return 0, nil, fmt.Errorf("id %s doesn't exist", id)
	}

	var edges []string
	for _, e := range l.Edges {
		if e[0] == id {
			edges = append(edges, e[1])
		}
	}

	return val, edges, nil
}
