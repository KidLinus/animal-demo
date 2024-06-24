package api

import (
	"fmt"

	"animal.dev/animal/internal/model"
)

type AnimalGet struct{ ID int }

func (app *App) AnimalGet(ctx Context, in AnimalGet) (*model.Animal, error) {
	itr, err := app.DB.AnimalGet(ctx, DatabaseAnimalGet{ID: []int{in.ID}, Limit: ptr(1)})
	if err != nil {
		return nil, err
	}
	res, err := iterOne(itr)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("animal not found")
	}
	return res, nil
}

type AnimalList struct {
	ID      []int
	Dataset *int
	Removed *bool
	Limit   *int
}

func (app *App) AnimalList(ctx Context, in AnimalList) ([]model.Animal, error) {
	itr, err := app.DB.AnimalGet(ctx, DatabaseAnimalGet{ID: in.ID, Dataset: in.Dataset, Removed: in.Removed, Limit: in.Limit})
	if err != nil {
		return nil, err
	}
	return iterAll(itr)
}

type AnimalCreate struct {
	ID      int
	Dataset int
	Name    string
}

func (app *App) AnimalCreate(ctx Context, in AnimalCreate) (*model.Animal, error) {
	mdl := model.Animal{ID: in.ID, Dataset: in.Dataset, Name: in.Name}
	if err := app.DB.AnimalCreate(ctx, mdl); err != nil {
		return nil, err
	}
	return &mdl, nil
}

type AnimalUpdate struct {
	ID      int
	Removed *bool
	Name    *string
}

func (app *App) AnimalUpdate(ctx Context, in AnimalUpdate) (*model.Animal, error) {
	update := DatabaseAnimalUpdate{Removed: in.Removed, Name: in.Name}
	res, err := app.DB.AnimalUpdate(ctx, in.ID, update)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("animal not found")
	}
	return res, nil
}

type AnimalMultipleFamily struct {
	ID       []int
	Distance int
}

func (app *App) AnimalMultipleFamily(ctx Context, in AnimalMultipleFamily) ([]model.Animal, error) {
	res := []model.Animal{}
	search, found := in.ID, map[int]struct{}{}
	for i := 0; i < in.Distance; i++ {
		if len(search) == 0 {
			break
		}
		itr, err := app.DB.AnimalGet(ctx, DatabaseAnimalGet{ID: search})
		if err != nil {
			return nil, err
		}
		search = nil
		animals, err := iterAll(itr)
		if err != nil {
			return nil, err
		}
		for _, animal := range animals {
			if _, ok := found[animal.ID]; ok {
				continue
			}
			found[animal.ID] = struct{}{}
			res = append(res, animal)
			for _, par := range animal.Parents {
				if _, ok := found[par]; ok {
					continue
				}
				search = append(search, par)
			}
		}
	}
	return res, nil
}

type AnimalFamily struct {
	ID       int
	Distance int
}

type AnimalFamilyResponse struct {
	model.Animal
	Family map[int]*model.Animal `json:"family"`
}

func (app *App) AnimalFamily(ctx Context, in AnimalFamily) (*AnimalFamilyResponse, error) {
	itr, err := app.DB.AnimalGet(ctx, DatabaseAnimalGet{ID: []int{in.ID}, Limit: ptr(1)})
	if err != nil {
		return nil, err
	}
	animal, err := iterOne(itr)
	if err != nil {
		return nil, err
	}
	if animal == nil {
		return nil, fmt.Errorf("animal not found")
	}
	family := map[int]*model.Animal{}
	// Get parents
	search := []int{}
	for _, id := range animal.Parents {
		search = append(search, id)
		family[id] = nil
	}
	for i := 0; i < in.Distance; i++ {
		if len(search) == 0 {
			break
		}
		itr, err := app.DB.AnimalGet(ctx, DatabaseAnimalGet{ID: search})
		if err != nil {
			return nil, err
		}
		search = nil
		parents, err := iterAll(itr)
		if err != nil {
			return nil, err
		}
		for idx, parent := range parents {
			family[parent.ID] = &parents[idx]
			for _, id := range parent.Parents {
				if _, ok := family[id]; !ok {
					search = append(search, id)
					family[id] = nil
				}
			}
		}
	}
	// Done
	return &AnimalFamilyResponse{Animal: *animal, Family: family}, nil
}
