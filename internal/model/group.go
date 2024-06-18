package model

import (
	"fmt"
	"io"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type Group struct {
	Animals map[int]Animal
}

func NewGroup() *Group { return &Group{Animals: map[int]Animal{}} }

func (group *Group) Validate() error {
	return nil
}

func (group *Group) Merge(other *Group) {
	for _, v := range other.Animals {
		group.AnimalAdd(v)
	}
}

func (group *Group) AnimalAdd(animal Animal) bool {
	_, ok := group.Animals[animal.ID]
	group.Animals[animal.ID] = animal
	return ok
}

func (group *Group) FilterFamily(id int, generations int) (*Group, error) {
	n := NewGroup()
	target, ok := group.Animals[id]
	if !ok {
		return nil, fmt.Errorf("Animal %d not found", id)
	}
	n.AnimalAdd(target)
	var search []int
	for _, id := range target.Parents {
		search = append(search, id)
	}
	for i := 0; i < generations-1; i++ {
		if len(search) == 0 {
			break
		}
		var found []int
		for id, v := range group.Animals {
			if sliceContains(search, id) {
				if !n.AnimalAdd(v) {
					for _, i := range v.Parents {
						found = append(found, i)
					}
				}
			}
		}
		search = found
	}
	return n, nil
}

func (group *Group) GraphDOT(out io.Writer) error {
	g := graph.New(graph.IntHash, graph.Directed())
	var edges []struct {
		from   int
		to     int
		gender Gender
	}
	for _, animal := range group.Animals {
		g.AddVertex(animal.ID, graph.VertexAttributes(map[string]string{
			"label":  fmt.Sprintf("%d - %s", animal.ID, animal.Name),
			"gender": animal.Gender.String(),
		}))
		for gender, id := range animal.Parents {
			edges = append(edges, struct {
				from   int
				to     int
				gender Gender
			}{from: id, to: animal.ID, gender: gender})
		}
	}
	for _, edge := range edges {
		c := "red"
		if edge.gender == Male {
			c = "blue"
		}
		g.AddEdge(edge.from, edge.to, graph.EdgeAttributes(map[string]string{
			"color": c,
		}))
	}
	return draw.DOT(g, os.Stdout)
}

func (group *Group) InbreedCoefficient(a, b uint64) (float64, error) {
	return 0, nil
}
