// nolint: gochecknoinits
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"alluvial/metric"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(metric.EthBalanceCounter)
}

func main() {
	var (
		logOpts = &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug, ReplaceAttr: nil}
		log     = slog.New(slog.NewTextHandler(os.Stdout, logOpts))
		addr    = "127.0.0.1:9000"
	)

	f := flag.NewFlagSet("alluvial", flag.ExitOnError)
	f.StringVar(&addr, "address", addr, "")
	_ = f.Parse(os.Args[1:])

	log.Info("server start")

	cfg, err := newConfig(addr)
	if err != nil {
		log.Error(err.Error())

		os.Exit(1)
	}

	appOpts := newOptions(cfg, log)

	app, cleanup, err := newApp(context.Background(), appOpts)
	if err != nil {
		log.Error(err.Error())

		os.Exit(1)
	}

	defer cleanup()

	run(context.Background(), app)

	log.Info("server stopped")
}

func run(ctx context.Context, app *App) {
	log := app.log

	srv, serverError := newServer(app)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-shutdown:
		log.Info("os", "signal", sig)
	case err := <-serverError:
		log.Error(err.Error())
	}

	log.Info("server shutdown")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}
}
