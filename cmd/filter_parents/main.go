package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"animal.dev/animal/internal/model"
)

var selectedID = flag.Uint64("id", 9999999, "animal to extract")
var depthMax = flag.Int("depth", 5, "Maximum parent depth")

func init() { flag.Parse() }

func main() {
	db := NewDatabase()
	dec := json.NewDecoder(os.Stdin)
	log.Println("reading")
	var idLast uint64
	for {
		var mdl model.Animal
		if err := dec.Decode(&mdl); err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		if err := db.AnimalAdd(mdl); err != nil {
			log.Fatalln("db: adding animal failed", err)
		}
		idLast = mdl.ID
	}
	log.Println("reading done", len(db.Animals))
	log.Println("building")
	db.Build()
	log.Println("building done")
	if *selectedID == 9999999 {
		selectedID = &idLast
	}
	animal := db.Animals[*selectedID]
	log.Println("selecting animal", animal.Animal)
	log.Println("getting parents")
	parents := animal.Parents()
	log.Println("getting parents done", len(parents))
	log.Println("output")
	enc := json.NewEncoder(os.Stdout)
	for _, animal := range parents {
		enc.Encode(animal.Animal)
	}
	log.Println("done")
}

func NewDatabase() *Database { return &Database{Animals: map[uint64]*Animal{}} }

type Database struct {
	Animals map[uint64]*Animal
}

type Animal struct {
	model.Animal
	db      *Database
	parents []*Animal
}

func (db *Database) AnimalAdd(model model.Animal) error {
	if _, ok := db.Animals[model.ID]; ok {
		return fmt.Errorf("animal %d already exists", model.ID)
	}
	db.Animals[model.ID] = &Animal{Animal: model, db: db}
	return nil
}

func (db *Database) Build() error {
	// Find parents
	for _, animal := range db.Animals {
		animal.parents = nil
		for _, parentID := range animal.Animal.Parents {
			if parent, ok := db.Animals[parentID]; ok {
				animal.parents = append(animal.parents, parent)
			}
		}
	}
	// Done
	return nil
}

func (animal *Animal) Parents() []*Animal {
	m := map[uint64]*Animal{animal.ID: animal}
	animal.parentsScan(0, m)
	arr := make([]*Animal, 0, len(m))
	for _, v := range m {
		arr = append(arr, v)
	}
	return arr
}

func (animal *Animal) parentsScan(depth int, m map[uint64]*Animal) {
	for _, parent := range animal.parents {
		if _, ok := m[parent.ID]; !ok {
			m[parent.ID] = parent
			if depth < *depthMax {
				parent.parentsScan(depth+1, m)
			}
		}
	}
}
