package main

import (
	"log/slog"
	"net/http"
)

type Options struct {
	Config     config
	httpClient *http.Client
	Log        *slog.Logger
}

func newOptions(c config, l *slog.Logger) Options {
	return Options{
		Config:     c,
		httpClient: &http.Client{Timeout: c.GatewayTimeout},
		Log:        l,
	}
}
