package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"animal.dev/animal/internal/model"
	"github.com/goombaio/namegenerator"
)

func init() { flag.Parse() }

var ancestors = flag.Uint64("ancestors", 1000, "How many original animals starts the simulation")
var years = flag.Uint64("years", 50, "How many years to simulate")

func main() {
	namegen := namegenerator.NewNameGenerator(time.Now().Unix())
	// Ancestors
	var idNext uint64
	animals := map[uint64]*model.Animal{}
	animalsAlive := map[uint64]*model.Animal{}
	for i := uint64(0); i < *ancestors; i++ {
		id := idNext
		idNext++
		animal := &model.Animal{ID: id, Born: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Name: namegen.Generate(), Gender: model.Gender(rand.Intn(2)), Parents: map[model.Gender]uint64{}}
		animals[animal.ID], animalsAlive[animal.ID] = animal, animal
	}
	dateStart := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.AddDate(int(*years), 0, 0)
	date := dateStart
	for {
		if !date.Before(dateEnd) {
			break
		}
		// Check population
		potentialMoms, potentialDads := []uint64{}, []uint64{}
		for _, animal := range animalsAlive {
			age := date.Sub(animal.Born)
			// Check if animal dies
			if rand.Intn(10) == 0 || age.Hours()/(24*365) > 10 {
				d := date
				animal.Died = &d
				delete(animalsAlive, animal.ID)
				continue
			}
			// Add to breeding pool
			if age > time.Hour*24*365 && age < time.Hour*24*365*8 {
				if animal.Gender == model.Female {
					potentialMoms = append(potentialMoms, animal.ID)
				} else {
					potentialDads = append(potentialDads, animal.ID)
				}
			}
		}
		// Breeding
		for _, mom := range potentialMoms {
			if len(potentialDads) > 0 && rand.Intn(10) == 0 {
				dad := potentialDads[rand.Intn(len(potentialDads))]
				kids := 2 + rand.Intn(7)
				for i := 0; i < kids; i++ {
					id := idNext
					idNext++
					animal := &model.Animal{ID: id, Born: date, Name: namegen.Generate(), Gender: model.Gender(rand.Intn(2)), Parents: map[model.Gender]uint64{model.Female: mom, model.Male: dad}}
					animals[animal.ID], animalsAlive[animal.ID] = animal, animal
				}
			}
		}
		date = date.AddDate(0, 6, 0)
	}
	log.Println("Done", len(animalsAlive), "/", len(animals))
	enc := json.NewEncoder(os.Stdout)
	for _, animal := range animals {
		enc.Encode(animal)
	}
}
