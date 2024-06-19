package api

import (
	"fmt"

	"animal.dev/animal/internal/model"
)

type DatasetGet struct{ ID int }

func (app *App) DatasetGet(ctx Context, in DatasetGet) (*model.Dataset, error) {
	itr, err := app.DB.DatasetGet(ctx, DatabaseDatasetGet{ID: []int{in.ID}, Limit: ptr(1)})
	if err != nil {
		return nil, err
	}
	res, err := iterOne(itr)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("dataset not found")
	}
	return res, nil
}

type DatasetList struct {
	ID      []int
	Removed *bool
	Limit   *int
}

func (app *App) DatasetList(ctx Context, in DatasetList) ([]model.Dataset, error) {
	itr, err := app.DB.DatasetGet(ctx, DatabaseDatasetGet{ID: in.ID, Removed: in.Removed, Limit: in.Limit})
	if err != nil {
		return nil, err
	}
	return iterAll(itr)
}

type DatasetCreate struct {
	ID   int
	Name string
}

func (app *App) DatasetCreate(ctx Context, in DatasetCreate) (*model.Dataset, error) {
	mdl := model.Dataset{ID: in.ID, Name: in.Name}
	if err := app.DB.DatasetCreate(ctx, mdl); err != nil {
		return nil, err
	}
	return &mdl, nil
}

type DatasetUpdate struct {
	ID      int
	Removed *bool
	Name    *string
}

func (app *App) DatasetUpdate(ctx Context, in DatasetUpdate) (*model.Dataset, error) {
	update := DatabaseDatasetUpdate{Removed: in.Removed, Name: in.Name}
	res, err := app.DB.DatasetUpdate(ctx, in.ID, update)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("dataset not found")
	}
	return res, nil
}
