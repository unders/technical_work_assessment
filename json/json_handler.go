package json

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/uptrace/bunrouter"
)

func HandlerNotFound(log *slog.Logger) func(http.ResponseWriter, bunrouter.Request) error {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		log.Warn(fmt.Sprintf("%s %s 404 Not Found", r.Method, r.URL.RequestURI()))

		return writeError(w, http.StatusNotFound, `{ "message": "Not Found" }`)
	}
}

func HandlerMethodNotAllowed(log *slog.Logger) func(http.ResponseWriter, bunrouter.Request) error {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		log.Warn(fmt.Sprintf("%s %s 405 Method Not Allowed", r.Method, r.URL.RequestURI()))

		return writeError(w, http.StatusMethodNotAllowed, `{ "message": "Method Not Allowed" }`)
	}
}
