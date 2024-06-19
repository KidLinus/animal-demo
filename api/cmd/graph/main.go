package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"animal.dev/animal/internal/model"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

// cat database.json | go run cmd/draw/main.go > graph.gv
// dot -Tsvg -O graph.gv

func main() {
	g := graph.New(graph.IntHash, graph.Directed())
	dec := json.NewDecoder(os.Stdin)
	log.Println("reading data")
	edges := []edge{}
	for {
		var animal model.Animal
		if err := dec.Decode(&animal); err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		gender := "male"
		if animal.Gender == model.Female {
			gender = "female"
		}
		if err := g.AddVertex(int(animal.ID), graph.VertexAttribute("gender", gender)); err != nil {
			log.Fatalln(err)
		}
		for gender, parentID := range animal.Parents {
			parent := "dad"
			if gender == model.Female {
				parent = "mom"
			}
			edges = append(edges, edge{parent: parentID, child: animal.ID, title: parent})
		}
	}
	log.Println("putting in edges", len(edges))
	for _, edge := range edges {
		if err := g.AddEdge(int(edge.parent), int(edge.child), graph.EdgeAttribute("parent", edge.title)); err != nil {
			//log.Fatalln(err)
		}
	}
	log.Println("building graph")
	if err := draw.DOT(g, os.Stdout); err != nil {
		log.Fatal(err)
	}
	log.Println("done")
}

type edge struct {
	parent uint64
	child  uint64
	title  string
}
