package main

import (
	"net/http"

	"alluvial/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/bunrouter"
)

func (app *App) handler() http.Handler {
	router := bunrouter.New(
		bunrouter.Use(),
		bunrouter.WithNotFoundHandler(json.HandlerNotFound(app.log)),
		bunrouter.WithMethodNotAllowedHandler(json.HandlerMethodNotAllowed(app.log)),
	).Compat()

	router.GET("/healthz", healthz(app).serveText)
	router.GET("/metrics", promhttp.Handler().ServeHTTP)
	router.GET("/eth/balance/:address", ethBalance(app).serveJSON)

	return router
}
