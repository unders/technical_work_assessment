package main

import (
	"context"
	"fmt"
	"log/slog"

	"alluvial/ethereum"
	"alluvial/json"
)

type App struct {
	address  string
	ethereum *ethereum.LoadBalancer
	log      *slog.Logger
	json     json.Writer
}

func newApp(ctx context.Context, opts Options) (*App, func(), error) {
	lbOpts := ethereum.Options{
		Gateways:   opts.Config.Gateways,
		HTTPClient: opts.httpClient,
		Log:        opts.Log,
		MinClients: opts.Config.GatewaysMinimumClients,
	}

	lb, err := ethereum.NewLoadBalancer(ctx, lbOpts)
	if err != nil {
		return &App{}, func() {}, fmt.Errorf("loadbalancer setup error: %w", err)
	}

	app := &App{
		address:  opts.Config.Address,
		json:     json.NewWriter(opts.Log),
		ethereum: lb,
		log:      opts.Log,
	}

	cleanup := func() {
		lb.Close()
	}

	return app, cleanup, nil
}
