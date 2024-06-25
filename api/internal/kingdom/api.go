package kingdom

import (
	"context"
	"fmt"
	"io"
)

type API struct {
	DB APIDatabase
}

type APIDatabase interface {
	Get(context.Context, string) (*Animal, error)
	List(context.Context, AnimalFilter) (AnimalIterator, error)
}

type APIAnimalSearchInput struct{ Query string }

type APIAnimalSearchResult struct{ Items []Animal }

func (api *API) AnimalSearch(ctx context.Context, in APIAnimalSearchInput) (*APIAnimalSearchResult, error) {
	itr, err := api.DB.List(ctx, AnimalFilter{Query: &in.Query, Limit: ptr(10)})
	if err != nil {
		return nil, err
	}
	items := []Animal{}
	for {
		var item Animal
		if err := itr.Next(&item); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		items = append(items, item)
	}
	return &APIAnimalSearchResult{Items: items}, nil
}

func (api *API) AnimalGet(ctx context.Context, id string) (*Animal, error) {
	v, err := api.DB.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, fmt.Errorf("animal %s not found", id)
	}
	return v, nil
}

type APIAnimalGetParentsInput struct {
	ID          string
	Depth       int
	FillUnknown bool
}

func (api *API) AnimalGetParents(ctx context.Context, in APIAnimalGetParentsInput) (*Group, error) {
	animal, err := api.AnimalGet(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	animals, search := map[string]*Animal{animal.ID: animal}, []string{}
	for _, id := range animal.Parents {
		if _, exists := animals[id]; !exists {
			search = append(search, id)
			animals[id] = nil
		}
	}
	for i := 0; i < in.Depth; i++ {
		itr, err := api.DB.List(ctx, AnimalFilter{IDs: search})
		if err != nil {
			return nil, err
		}
		for {
			var animal Animal
			if err := itr.Next(&animal); err != nil {
				if err != io.EOF {
					return nil, err
				}
				break
			}
			animals[animal.ID] = &animal
			for _, id := range animal.Parents {
				if _, exists := animals[id]; !exists {
					search = append(search, id)
					animals[id] = nil
				}
			}
		}
	}
	group := NewGroup()
	for _, animal := range animals {
		if animal != nil {
			group.Add(*animal)
		}
	}
	return group, nil
}

type APIAnimalGetParentTreeInput struct {
	ID          string
	Depth       int
	FillUnknown bool
}

func (api *API) AnimalGetParentTree(ctx context.Context, in APIAnimalGetParentTreeInput) (*Tree, error) {
	group, err := api.AnimalGetParents(ctx, APIAnimalGetParentsInput{ID: in.ID, Depth: in.Depth, FillUnknown: in.FillUnknown})
	tree, err := group.TreeAnimalParents(in.ID, in.Depth, in.FillUnknown)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

type APIAnimalGetCOIInput struct {
	ID    string
	Depth int
}

func (api *API) AnimalGetCOI(ctx context.Context, in APIAnimalGetCOIInput) (*AnimalInbreedingCoefficient, error) {
	group, err := api.AnimalGetParents(ctx, APIAnimalGetParentsInput{ID: in.ID, Depth: in.Depth})
	if err != nil {
		return nil, err
	}
	return group.AnimalInbreedingCoefficient(in.ID, in.Depth)
}
