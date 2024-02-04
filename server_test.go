package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func newTestServer(t *testing.T, client *http.Client) (*httptest.Server, func()) {
	t.Helper()

	cfg, err := newConfig("address is not used when running tests")
	if err != nil {
		t.Fatal(err)
	}

	url, ok := os.LookupEnv("ALCHEMY")
	if !ok {
		t.Fatal("env variable ALCHEMY must be set")
	}

	// we only use one endpoint in these tests in order to
	// simplify testing and recording of the requests.
	cfg.Gateways = map[string]string{
		"Alchemy": url,
	}

	cfg.GatewaysMinimumClients = 1

	appOpts := newOptions(cfg, slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, AddSource: true})))

	if client != nil {
		appOpts.httpClient = client
	}

	app, cleanup, err := newApp(context.Background(), appOpts)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(app.handler())

	return ts, func() {
		cleanup()
		ts.Close()
	}
}
