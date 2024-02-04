package main

import (
	"context"
	"errors"
)

func (app *App) Healthz(ctx context.Context) error {
	if ok := app.ethereum.HealthCheck(ctx); !ok {
		return errors.New("all ethereum backends is down")
	}

	return nil
}
