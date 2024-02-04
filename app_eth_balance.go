package main

import (
	"context"
	"math/big"

	"alluvial/ethereum"
)

func (app *App) EthBalance(ctx context.Context, req *requestEthBalance, resp *responseEthBalance) error {
	var balance *big.Int

	f := func(c ethereum.Client) error {
		var err error

		balance, err = c.BalanceAtLatestBlock(ctx, req.Address)

		return err
	}

	if err := app.ethereum.Call(f); err != nil {
		return err
	}

	resp.Balance = balance.String()

	return nil
}
