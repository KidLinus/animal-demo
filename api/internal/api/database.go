package api

import (
	"context"
	"io"

	"animal.dev/animal/internal/model"
)

type Database interface {
	DatasetCreate(context.Context, model.Dataset) error
	DatasetGet(context.Context, DatabaseDatasetGet) (DatabaseDatasetIterator, error)
	DatasetUpdate(context.Context, int, DatabaseDatasetUpdate) (*model.Dataset, error)

	AnimalCreate(context.Context, model.Animal) error
	AnimalGet(context.Context, DatabaseAnimalGet) (DatabaseAnimalIterator, error)
	AnimalUpdate(context.Context, int, DatabaseAnimalUpdate) (*model.Animal, error)
}

type DatabaseDatasetGet struct {
	ID      []int
	Removed *bool
	Limit   *int
}

type DatabaseDatasetUpdate struct {
	Removed *bool
	Name    *string
}

type DatabaseDatasetIterator interface{ Next(*model.Dataset) error }

type DatabaseAnimalGet struct {
	ID      []int
	Parent  []int
	Dataset *int
	Removed *bool
	Limit   *int
}

type DatabaseAnimalUpdate struct {
	Removed *bool
	Name    *string
}

type DatabaseAnimalIterator interface{ Next(*model.Animal) error }

func iterAll[V any](itr interface{ Next(*V) error }) ([]V, error) {
	res := []V{}
	for {
		var mdl V
		if err := itr.Next(&mdl); err != nil {
			if err == io.EOF {
				return res, nil
			}
			return nil, err
		}
		res = append(res, mdl)
	}
}

func iterOne[V any](itr interface{ Next(*V) error }) (*V, error) {
	var mdl V
	if err := itr.Next(&mdl); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	return &mdl, nil
}
