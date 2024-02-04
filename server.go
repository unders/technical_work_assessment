package main

import (
	"log/slog"
	"net/http"
	"time"
)

func newServer(app *App) (*http.Server, <-chan error) {
	log := app.log

	srv := &http.Server{
		Addr:              app.address,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       6 * time.Second,
		WriteTimeout:      9 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes, // Default is 1MB.
		Handler:           app.handler(),
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelDebug),
	}

	errChan := make(chan error, 1)

	go func() {
		log.Info("listens on", "address", app.address)

		errChan <- srv.ListenAndServe()
	}()

	return srv, errChan
}
