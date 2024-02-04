// nolint: goconst
package main

import (
	"net/http"
	"testing"
	"time"

	"alluvial/test"
)

func TestHandler_EthBalance_WhenValidAddress(t *testing.T) {
	tr := test.TransportRecorder(false, false, "handler_ethBalance_when_valid_address")

	c := http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	ts, cleanup := newTestServer(t, &c)
	defer cleanup()

	resp := &responseEthBalance{}

	test.GetJSON(t, ts.URL+"/eth/balance/0xfe3b557e8fb62b89f4916b721be55ceb828dbd73", http.StatusOK, resp)

	wantBalance := "14058"
	if wantBalance != resp.Balance {
		t.Errorf("\nWant: %s\n Got: %s\n", wantBalance, resp.Balance)
	}
}

func TestHandler_EthBalance_WhenInvalidAddress(t *testing.T) {
	tr := test.TransportRecorder(false, false, "handler_ethBalance_when_invalid_address")

	c := http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	ts, cleanup := newTestServer(t, &c)
	defer cleanup()

	resp := struct {
		Message string `json:"message"`
	}{}

	test.GetJSON(t, ts.URL+"/eth/balance/0xfe3b557e8fb62b89f4916b721be55ceb828dbd73xx", http.StatusBadRequest, &resp)

	wantMessage := "invalid ethereum address"
	if wantMessage != resp.Message {
		t.Errorf("\nWant: %s\n Got: %s\n", wantMessage, resp.Message)
	}
}

func TestHandler_EthBalance_WhenContractAddress(t *testing.T) {
	tr := test.TransportRecorder(false, false, "handler_ethBalance_when_contract_Address")

	c := http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	ts, cleanup := newTestServer(t, &c)
	defer cleanup()

	resp := struct {
		Message string `json:"message"`
	}{}

	test.GetJSON(t, ts.URL+"/eth/balance/0xe41d2489571d322189246dafa5ebde1f4699f498", http.StatusBadRequest, &resp)

	wantMessage := "ethereum contract address is not supported"
	if wantMessage != resp.Message {
		t.Errorf("\nWant: %s\n Got: %s\n", wantMessage, resp.Message)
	}
}

func TestHandler_EthBalance_WhenInternalError(t *testing.T) {
	tr := test.TransportRecorder(false, false, "handler_ethBalance_when_internal_error")

	c := http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	ts, cleanup := newTestServer(t, &c)
	defer cleanup()

	resp := struct {
		Message string `json:"message"`
	}{}

	test.GetJSON(t, ts.URL+"/eth/balance/0xe41d2489571d322189246dafa5ebde1f4699f498", http.StatusInternalServerError, &resp)

	wantMessage := "Internal Server Error"
	if wantMessage != resp.Message {
		t.Errorf("\nWant: %s\n Got: %s\n", wantMessage, resp.Message)
	}
}

func TestHandler_EthBalance_WhenInvalidJSON(t *testing.T) {
	tr := test.TransportRecorder(false, false, "handler_ethBalance_when_invalid_json")

	c := http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	ts, cleanup := newTestServer(t, &c)
	defer cleanup()

	resp := struct {
		Message string `json:"message"`
	}{}

	test.GetJSON(t, ts.URL+"/eth/balance/0xe41d2489571d322189246dafa5ebde1f4699f498", http.StatusInternalServerError, &resp)

	wantMessage := "Internal Server Error"
	if wantMessage != resp.Message {
		t.Errorf("\nWant: %s\n Got: %s\n", wantMessage, resp.Message)
	}
}
