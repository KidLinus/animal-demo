package main

import (
	"log"
	"os"

	"animal.dev/animal/internal/model"
)

func main() {
	all := model.NewGroup()
	all.AnimalAdd(model.Animal{ID: 1, Name: "Kalle"})
	all.AnimalAdd(model.Animal{ID: 2, Name: "Kallina"})
	all.AnimalAdd(model.Animal{ID: 3, Name: "Stina", Parents: map[model.Gender]int{model.Male: 1, model.Female: 2}})
	all.AnimalAdd(model.Animal{ID: 4, Name: "Erik", Parents: map[model.Gender]int{model.Male: 1, model.Female: 2}})
	all.AnimalAdd(model.Animal{ID: 5, Name: "Sven", Parents: map[model.Gender]int{model.Male: 1, model.Female: 2}})
	all.AnimalAdd(model.Animal{ID: 6, Name: "Eva", Parents: map[model.Gender]int{model.Male: 4, model.Female: 7}})
	all.AnimalAdd(model.Animal{ID: 7, Name: "Okina"})
	all.AnimalAdd(model.Animal{ID: 8, Name: "Göran", Parents: map[model.Gender]int{model.Male: 4, model.Female: 3}})
	all.AnimalAdd(model.Animal{ID: 9, Name: "Lisa", Parents: map[model.Gender]int{model.Male: 8, model.Female: 2}})
	all.AnimalAdd(model.Animal{ID: 10, Name: "Greta", Parents: map[model.Gender]int{model.Male: 8, model.Female: 9}})
	if err := all.Validate(); err != nil {
		log.Fatal(err)
	}
	family1, err := all.FilterFamily(6, 4)
	if err != nil {
		log.Fatal(err)
	}
	family2, err := all.FilterFamily(10, 4)
	if err != nil {
		log.Fatal(err)
	}
	family1.Merge(family2)
	family1.GraphDOT(os.Stdout)
}
