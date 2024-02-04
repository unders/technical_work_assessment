package json

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Writer struct {
	log *slog.Logger
}

func NewWriter(log *slog.Logger) Writer {
	return Writer{log: log}
}

func (w Writer) Write(rw http.ResponseWriter, r *http.Request, code int, headers map[string]string, body []byte) error {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	for k, v := range headers {
		rw.Header().Set(k, v)
	}

	rw.WriteHeader(code)

	if _, err := io.Copy(rw, bytes.NewReader(body)); err != nil {
		w.log.Error(err.Error(), "req", fmt.Sprintf("%s %s", r.Method, r.URL.RequestURI()))

		return errors.New("http write error")
	}

	if _, err := io.Copy(rw, bytes.NewBufferString("\n")); err != nil {
		w.log.Warn(err.Error())
	}

	return nil
}

func (w Writer) WriteInternalError(rw http.ResponseWriter) {
	if err := writeError(rw, http.StatusInternalServerError, `{ "message": "Internal Server Error" }`); err != nil {
		w.log.Error(err.Error())
	}
}

func (w Writer) WriteBadRequest(rw http.ResponseWriter, errMessage string) {
	body := fmt.Sprintf(`{ "message" : %q }`, errMessage)

	if err := writeError(rw, http.StatusBadRequest, body); err != nil {
		w.log.Error(err.Error())
	}
}

func writeError(w http.ResponseWriter, code int, body string) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)

	if _, err := fmt.Fprintln(w, body); err != nil {
		return fmt.Errorf("http write error: %w", err)
	}

	return nil
}
