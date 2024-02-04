package main

import (
	"context"
	stdjson "encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"alluvial/ethereum"
	"alluvial/json"
	"alluvial/metric"
	"github.com/ethereum/go-ethereum/common"
	"github.com/uptrace/bunrouter"
)

func ethBalance(app *App) handlerEthBalance {
	return handlerEthBalance{log: app.log, json: app.json, app: app}
}

type ethBalanceAppInterface interface {
	EthAddressIsContract(ctx context.Context, address common.Address) (bool, error)
	EthBalance(ctx context.Context, req *requestEthBalance, resp *responseEthBalance) error
}

type handlerEthBalance struct {
	log  *slog.Logger
	json json.Writer
	app  ethBalanceAppInterface
}

type requestEthBalance struct {
	Address common.Address
}

type responseEthBalance struct {
	Balance string `json:"balance"`
}

func (h handlerEthBalance) serveJSON(w http.ResponseWriter, r *http.Request) {
	metric.EthBalanceCounter.Inc()

	var (
		ctx  = r.Context()
		req  = &requestEthBalance{}
		resp = &responseEthBalance{}
	)

	if err := h.requestRead(ctx, w, req, r); err != nil {
		return
	}

	isContract, err := h.app.EthAddressIsContract(ctx, req.Address)
	if err != nil {
		h.log.Error(err.Error(), "req", fmt.Sprintf("%s %s 500 Internal Server Error", r.Method, r.URL.RequestURI()))

		h.json.WriteInternalError(w)

		return
	}

	if isContract {
		const msg = "ethereum contract address is not supported"

		h.log.Warn(msg, "req", fmt.Sprintf("%s %s 400 Bad Request", r.Method, r.URL.RequestURI()))

		h.json.WriteBadRequest(w, msg)

		return
	}

	if err = h.app.EthBalance(ctx, req, resp); err != nil {
		h.log.Error(err.Error(), "req", fmt.Sprintf("%s %s 500 Internal Server Error", r.Method, r.URL.RequestURI()))

		h.json.WriteInternalError(w)

		return
	}

	body, err := stdjson.Marshal(resp)
	if err != nil {
		h.log.Error(err.Error(), "req", fmt.Sprintf("%s %s 500 Internal Server Error", r.Method, r.URL.RequestURI()))

		h.json.WriteInternalError(w)

		return
	}

	if err = h.json.Write(w, r, http.StatusOK, nil, body); err != nil {
		return
	}

	h.log.Info(fmt.Sprintf("%s %s 200 OK", r.Method, r.URL.RequestURI()))
}

func (h handlerEthBalance) requestRead(ctx context.Context, w http.ResponseWriter, req *requestEthBalance, r *http.Request) error {
	params := bunrouter.ParamsFromContext(ctx)

	addr, found := params.Get("address")
	if !found {
		const msg = "missing ethereum address"

		h.log.Warn(msg, "req", fmt.Sprintf("%s %s 400 Bad Request", r.Method, r.URL.RequestURI()))

		h.json.WriteBadRequest(w, msg)

		return ErrAbort
	}

	msg, ethAddr, valid := ethereum.IsValidAddress(addr)
	if !valid {
		h.log.Warn(msg, "req", fmt.Sprintf("%s %s 400 Bad Request", r.Method, r.URL.RequestURI()))

		h.json.WriteBadRequest(w, msg)

		return ErrAbort
	}

	req.Address = ethAddr

	return nil
}
