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
