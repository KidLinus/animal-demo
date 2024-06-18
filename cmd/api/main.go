package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"animal.dev/animal/internal/api"
	"animal.dev/animal/internal/api/http"
	"animal.dev/animal/internal/model"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Addr         []string
	AllowOrigins []string
}

func init() { gin.SetMode(gin.ReleaseMode) }

func main() {
	var cfg Config
	cfg.Addr = []string{":8666"}
	cfg.AllowOrigins = []string{"*"}
	db := &memstore{datasets: map[int]model.Dataset{
		1: {ID: 1, Name: "dogs"},
		2: {ID: 2, Name: "cats", Removed: true},
	}}
	app, err := api.New(api.Config{DB: db})
	if err != nil {
		log.Fatal("api create fail", err)
	}
	srv, err := http.New(http.Config{Addr: cfg.Addr, AllowOrigins: cfg.AllowOrigins, API: app})
	if err != nil {
		log.Fatal("http create fail", err)
	}
	log.Println("running")
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case sig := <-ch:
			log.Println("shutdown signal:", sig)
			srv.Close()
			app.Close()
			log.Println("closed")
			return
		}
	}
}

type memstore struct {
	datasets map[int]model.Dataset
	animals  map[int]model.Animal
}

type memstoreDatasetItr struct {
	db  *memstore
	ids []int
	i   int
}

func (db *memstore) DatasetCreate(_ context.Context, mdl model.Dataset) error {
	if _, ok := db.datasets[mdl.ID]; ok {
		return fmt.Errorf("dataset %d exists", mdl.ID)
	}
	db.datasets[mdl.ID] = mdl
	return nil
}

func (db *memstore) DatasetGet(_ context.Context, filter api.DatabaseDatasetGet) (api.DatabaseDatasetIterator, error) {
	itr := &memstoreDatasetItr{db: db}
	for id, mdl := range db.datasets {
		if len(filter.ID) > 0 && !sliceContains(filter.ID, id) {
			continue
		}
		if filter.Removed != nil && mdl.Removed != *filter.Removed {
			continue
		}
		itr.ids = append(itr.ids, id)
		if filter.Limit != nil && len(itr.ids) >= *filter.Limit {
			break
		}
	}
	return itr, nil
}

func (db *memstore) DatasetUpdate(_ context.Context, id int, update api.DatabaseDatasetUpdate) (*model.Dataset, error) {
	mdl, ok := db.datasets[id]
	if !ok {
		return nil, nil
	}
	if update.Removed != nil {
		mdl.Removed = *update.Removed
	}
	if update.Name != nil {
		mdl.Name = *update.Name
	}
	db.datasets[id] = mdl
	return &mdl, nil
}

func (itr *memstoreDatasetItr) Next(v *model.Dataset) error {
	if len(itr.ids) > itr.i {
		*v = itr.db.datasets[itr.ids[itr.i]]
		itr.i++
		return nil
	}
	return io.EOF
}

func sliceContains[T comparable](s []T, v T) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}
