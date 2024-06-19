package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type Katt struct {
	KattID         int      `json:"katt_id"`
	MorID          *int     `json:"mor_id"`
	FarID          *int     `json:"far_id"`
	Fodelsedatum   string   `json:"fodelsedatum"`
	Farg           string   `json:"farg"`
	Forbund        string   `json:"forbund"`
	Stambok        string   `json:"stambok"`
	StamnamnPrefix string   `json:"stamnamn_prefix"`
	Stamnamn       string   `json:"stamnamn"`
	StamnamnSuffix string   `json:"stamnamn_suffix"`
	Regnr          string   `json:"regnr"`
	IDMarkningsnr  *int     `json:"id_markningsnr"`
	Namn           string   `json:"namn"`
	Ras            string   `json:"ras"`
	Genotyp        string   `json:"genotyp"`
	Kon            string   `json:"kon"`
	Kastrat        string   `json:"kastrat"`
	Kattstatus     string   `json:"kattstatus"`
	Isexport       string   `json:"isexport"`
	Isimport       string   `json:"isimport"`
	Regnummerlista []string `json:"regnummerlista"`
}

func main() {
	g := graph.New(graph.IntHash, graph.Directed())
	dec := json.NewDecoder(os.Stdin)
	katter := map[int]Katt{}
	edges := []edge{}
	k := []Katt{}
	dec.Decode(&k)
	for _, katt := range k {
		katter[katt.KattID] = katt
		g.AddVertex(katt.KattID, graph.VertexAttributes(map[string]string{"namn": katt.Namn}))
		if katt.MorID != nil {
			edges = append(edges, edge{parent: *katt.MorID, child: katt.KattID, title: "mor"})
		}
		if katt.FarID != nil {
			edges = append(edges, edge{parent: *katt.FarID, child: katt.KattID, title: "far"})
		}
	}
	log.Println("putting in edges", len(edges))
	for _, edge := range edges {
		if err := g.AddEdge(int(edge.parent), int(edge.child), graph.EdgeAttribute("parent", edge.title), graph.EdgeWeight(1)); err != nil {
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
	parent int
	child  int
	title  string
}
