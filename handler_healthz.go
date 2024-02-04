package main

import (
	"fmt"
	"net/http"
)

func healthz(app *App) handlerHealthz {
	return handlerHealthz{app: app}
}

type handlerHealthz struct{ app *App }

func (h handlerHealthz) serveText(w http.ResponseWriter, r *http.Request) {
	log := h.app.log

	if err := h.app.Healthz(r.Context()); err != nil {
		msg := fmt.Sprintf("service is down! error: %s", err)
		req := fmt.Sprintf("%s %s 500 Internal Server Error", r.Method, r.URL.RequestURI())

		log.Error(msg, "req", req)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(msg))
		_, _ = w.Write([]byte("\n"))

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))

	log.Info(fmt.Sprintf("%s %s 200 OK", r.Method, r.URL.RequestURI()))
}
