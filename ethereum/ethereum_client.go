package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	client *ethclient.Client
}

func Dial(ctx context.Context, url string, c *http.Client) (Client, error) {
	rpcClient, err := rpc.DialOptions(ctx, url, rpc.WithHTTPClient(c))
	if err != nil {
		return Client{}, fmt.Errorf("rpc dial error: %w", err)
	}

	return Client{client: ethclient.NewClient(rpcClient)}, nil
}

// BalanceAtLatestBlock the balance is taken from the latest known block for the supplied address account.
func (c Client) BalanceAtLatestBlock(ctx context.Context, addr common.Address) (*big.Int, error) {
	b, err := c.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return b, fmt.Errorf("balance at error: %w", err)
	}

	return b, nil
}

func (c Client) IsContract(ctx context.Context, addr common.Address) (bool, error) {
	bytecode, err := c.client.CodeAt(ctx, addr, nil)
	if err != nil {
		return false, fmt.Errorf("code at error: %w", err)
	}

	return len(bytecode) > 0, nil
}

func (c Client) ChainID(ctx context.Context) (*big.Int, error) {
	id, err := c.client.ChainID(ctx)
	if err != nil {
		return id, fmt.Errorf("chain id error: %w", err)
	}

	return id, nil
}

func (c Client) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
