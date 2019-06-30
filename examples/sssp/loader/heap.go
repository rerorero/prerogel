package loader

import (
	"github.com/rerorero/prerogel/plugin"
)

var (
	// ref. https://www.cse.cuhk.edu.hk/~taoyf/course/comp3506/tut/tut12.pdf
	vertices = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	edges    = []struct {
		src  string
		dest string
		dist uint32
	}{
		{"a", "b", 2},
		{"b", "d", 3},
		{"c", "a", 1},
		{"c", "b", 5},
		{"d", "c", 7},
		{"d", "e", 1},
		{"e", "f", 1},
		{"f", "c", 3},
		{"f", "h", 4},
		{"g", "d", 6},
		{"g", "i", 2},
		{"h", "g", 1},
		{"h", "i", 4},
	}
	// expected results
	// id	dist 	parent
	// a	0		nil
	// b	2		a
	// c	10		f
	// d	5		b
	// e	6		d
	// f	7		e
	// g	12		h
	// h	11		f
	// i	14		g
)

// HeapLoader loads vertex from heap
type HeapLoader struct{}

// LoadPartition loads vertices of partition
func (l *HeapLoader) LoadPartition(partitionID uint64, numOfPartitions uint64) ([]*Vert, error) {
	var vs []*Vert

	for _, id := range vertices {
		partition, err := plugin.HashPartition(plugin.VertexID(id), numOfPartitions)
		if err != nil {
			return nil, err
		}

		if partition == partitionID {
			v := &Vert{
				ID:        id,
				Outgoings: make(map[string]uint32),
			}

			for _, e := range edges {
				if e.src == id {
					v.Outgoings[e.dest] = e.dist
				}
			}

			vs = append(vs, v)
		}
	}

	return vs, nil
}
