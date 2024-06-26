package main

import (
	"context"
	"io"
	"strings"

	"animal.dev/animal/internal/kingdom"
)

var animals = map[string]*kingdom.Animal{
	"1":  {ID: "1", Name: "Fluffy", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{}},
	"2":  {ID: "2", Name: "Tipsy", Gender: kingdom.Female, Parents: map[kingdom.Gender]string{}},
	"3":  {ID: "3", Name: "Shiro", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{}},
	"4":  {ID: "4", Name: "Maja", Gender: kingdom.Female, Parents: map[kingdom.Gender]string{}},
	"5":  {ID: "5", Name: "Oreo", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{kingdom.Male: "1", kingdom.Female: "2"}},
	"6":  {ID: "6", Name: "Luna", Gender: kingdom.Female, Parents: map[kingdom.Gender]string{kingdom.Male: "3", kingdom.Female: "4"}},
	"7":  {ID: "7", Name: "Luffsen", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{kingdom.Male: "5", kingdom.Female: "6"}},
	"8":  {ID: "8", Name: "Missy", Gender: kingdom.Female, Parents: map[kingdom.Gender]string{kingdom.Male: "5", kingdom.Female: "7"}},
	"9":  {ID: "9", Name: "Lennart", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{kingdom.Male: "7", kingdom.Female: "8"}},
	"10": {ID: "10", Name: "Dobby", Gender: kingdom.Male, Parents: map[kingdom.Gender]string{kingdom.Male: "9", kingdom.Female: "8"}},
}

func main() {
	db := &db{animals: animals}
	api := kingdom.API{DB: db}
	gin := api.Gin([]string{"http://localhost:5173"})
	gin.Run(":8667")
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
		if filter.IDs != nil && !sliceContainsAny(filter.IDs, animal.ID) {
			continue
		}
		if filter.Query != nil {
			if !strings.Contains(animal.Name, *filter.Query) {
				continue
			}
		}
		itr.items = append(itr.items, *animal)
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

func sliceContainsAny[V comparable](slice []V, target ...V) bool {
	for _, v := range slice {
		for _, t := range target {
			if v == t {
				return true
			}
		}
	}
	return false
}
