package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"animal.dev/animal/internal/kingdom"
)

var animals = map[string]*kingdom.Animal{
	"1":  {ID: "1", Name: "1", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{}},
	"2":  {ID: "2", Name: "2", Gender: kingdom.Female, Parents: map[kingdom.Gender]string{}},
	"3":  {ID: "3", Name: "3", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{}},
	"4":  {ID: "4", Name: "4", Gender: kingdom.Female, Parents: map[kingdom.Gender]string{}},
	"5":  {ID: "5", Name: "5", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{kingdom.Male: "1", kingdom.Female: "2"}},
	"6":  {ID: "6", Name: "6", Gender: kingdom.Female, Parents: map[kingdom.Gender]string{kingdom.Male: "3", kingdom.Female: "4"}},
	"7":  {ID: "7", Name: "7", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{kingdom.Male: "5", kingdom.Female: "6"}},
	"8":  {ID: "8", Name: "8", Gender: kingdom.Female, Parents: map[kingdom.Gender]string{kingdom.Male: "5", kingdom.Female: "7"}},
	"9":  {ID: "9", Name: "9", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{kingdom.Male: "7", kingdom.Female: "8"}},
	"10": {ID: "10", Name: "10", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{kingdom.Male: "9", kingdom.Female: "8"}},
}

func main() {
	db := &db{animals: animals}
	api := kingdom.API{DB: db}
	coi, err := api.AnimalGetCOI(context.Background(), kingdom.APIAnimalGetCOIInput{ID: "9", Depth: 20})
	if err != nil {
		log.Fatalln("tree get fail", err)
	}
	js, _ := json.MarshalIndent(coi, "", "  ")
	log.Println(string(js))
}

type db struct {
	animals map[string]*kingdom.Animal
}

func (db *db) Get(_ context.Context, id string) (*kingdom.Animal, error) {
	return db.animals[id], nil
}

func (db *db) List(_ context.Context, filter kingdom.AnimalFilter) (kingdom.AnimalIterator, error) {
	itr := &dbItr{}
	for _, animal := range db.animals {
		for _, id := range filter.IDs {
			if animal.ID == id {
				itr.items = append(itr.items, *animal)
			}
		}
	}
	return itr, nil
}

type dbItr struct {
	index int
	items []kingdom.Animal
}

func (itr *dbItr) Next(v *kingdom.Animal) error {
	if itr.index > len(itr.items)-1 {
		return io.EOF
	}
	*v = itr.items[itr.index]
	itr.index++
	return nil
}
