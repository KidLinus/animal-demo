package http

import (
	"animal.dev/animal/internal/api"
	"animal.dev/animal/internal/model"
)

func (server *Server) animalList(ctx api.Context, in struct {
	ID      []int `form:"id" binding:"omitempty,min=1"`
	Dataset *int  `form:"dataset"`
	Removed *bool `form:"removed" default:"false"`
	Limit   *int  `form:"limit" default:"50" binding:"omitempty,min=1,max=100"`
}) ([]model.Animal, error) {
	return server.API.AnimalList(ctx, api.AnimalList(in))
}

func (server *Server) animalCreate(ctx api.Context, in struct {
	ID      int    `json:"id" binding:"required"`
	Dataset int    `json:"dataset" binding:"required"`
	Name    string `json:"name" binding:"required,min=1"`
}) (*model.Animal, error) {
	return server.API.AnimalCreate(ctx, api.AnimalCreate(in))
}

func (server *Server) animalGet(ctx api.Context, in struct {
	ID int `uri:"id" binding:"required"`
}) (*model.Animal, error) {
	return server.API.AnimalGet(ctx, api.AnimalGet(in))
}

func (server *Server) animalUpdate(ctx api.Context, in struct {
	ID      int     `uri:"id" binding:"required"`
	Removed *bool   `json:"removed"`
	Name    *string `json:"name" binding:"omitempty,min=1"`
}) (*model.Animal, error) {
	return server.API.AnimalUpdate(ctx, api.AnimalUpdate(in))
}

func (server *Server) animalRemove(ctx api.Context, in struct {
	ID int `uri:"id" binding:"required"`
}) (*model.Animal, error) {
	return server.API.AnimalUpdate(ctx, api.AnimalUpdate{ID: in.ID, Removed: ptr(true)})
}

func (server *Server) animalFamily(ctx api.Context, in struct {
	ID       int `uri:"id" binding:"required"`
	Distance int `form:"distance" default:"6" binding:"omitempty,max=10"`
}) (*api.AnimalFamilyResponse, error) {
	return server.API.AnimalFamily(ctx, api.AnimalFamily(in))
}

func (server *Server) animalMultipleFamily(ctx api.Context, in struct {
	ID       []int `form:"id" binding:"required,min=1,unique"`
	Distance int   `form:"distance" default:"6" binding:"omitempty,max=10"`
}) ([]model.Animal, error) {
	return server.API.AnimalMultipleFamily(ctx, api.AnimalMultipleFamily(in))
}
