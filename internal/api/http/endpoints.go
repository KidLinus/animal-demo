package http

import (
	"animal.dev/animal/internal/api"
	"animal.dev/animal/internal/model"
)

func (server *Server) datasetList(ctx api.Context, in struct {
	ID      []int `form:"id" binding:"omitempty,min=1"`
	Removed *bool `form:"removed" default:"false"`
	Limit   *int  `form:"limit" binding:"omitempty,min=1,max=100"`
}) ([]model.Dataset, error) {
	return server.API.DatasetList(ctx, api.DatasetList(in))
}

func (server *Server) datasetCreate(ctx api.Context, in struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required,min=1"`
}) (*model.Dataset, error) {
	return server.API.DatasetCreate(ctx, api.DatasetCreate(in))
}

func (server *Server) datasetGet(ctx api.Context, in struct {
	ID int `uri:"id" binding:"required"`
}) (*model.Dataset, error) {
	return server.API.DatasetGet(ctx, api.DatasetGet(in))
}

func (server *Server) datasetUpdate(ctx api.Context, in struct {
	ID      int     `uri:"id" binding:"required"`
	Removed *bool   `json:"removed"`
	Name    *string `json:"name" binding:"omitempty,min=1"`
}) (*model.Dataset, error) {
	return server.API.DatasetUpdate(ctx, api.DatasetUpdate(in))
}
