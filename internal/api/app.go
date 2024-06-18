package api

import (
	"context"
)

type App struct {
	Config
}

type Config struct{ DB Database }

func New(cfg Config) (*App, error) {
	app := &App{Config: cfg}
	return app, nil
}

type Context interface{ context.Context }

func (app *App) Close() {}
