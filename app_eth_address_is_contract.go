package main

import (
	"context"

	"alluvial/ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (app *App) EthAddressIsContract(ctx context.Context, address common.Address) (bool, error) {
	isContract := false

	f := func(c ethereum.Client) error {
		var err error

		isContract, err = c.IsContract(ctx, address)

		return err
	}

	err := app.ethereum.Call(f)

	return isContract, err
}
