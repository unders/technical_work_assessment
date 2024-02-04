package ethereum_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"alluvial/ethereum"
	"alluvial/test"
	"github.com/ethereum/go-ethereum/common"
)

func TestClient_BalanceAtLatestBlock(t *testing.T) {
	tr := test.TransportRecorder(false, false, "balance_at_latest_block_200_ok")

	c := newTestClient(t, "ALCHEMY", tr)
	defer c.Close()

	balance, err := c.BalanceAtLatestBlock(context.Background(), common.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73"))
	if err != nil {
		t.Fatalf("balance at latest block error: %s", err)
	}

	wantBalance := "14058"
	if wantBalance != balance.String() {
		t.Errorf("\nWant: %s\n Got: %s", wantBalance, balance.String())
	}
}

func TestClient_IsContract_WhenNotAContract(t *testing.T) {
	tr := test.TransportRecorder(false, false, "is_contract_200_ok_NotAContract")

	c := newTestClient(t, "ALCHEMY", tr)
	defer c.Close()

	isContract, err := c.IsContract(context.Background(), common.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73"))
	if err != nil {
		t.Errorf("\nWant: <nil>\n Got: %s", err)
	}

	if isContract {
		t.Errorf("\nWant: false\n Got: %t", isContract)
	}
}

func TestClient_IsContract_WhenContract(t *testing.T) {
	tr := test.TransportRecorder(false, false, "is_contract_200_ok_WhenContract")

	c := newTestClient(t, "ALCHEMY", tr)
	defer c.Close()

	isContract, err := c.IsContract(context.Background(), common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"))
	if err != nil {
		t.Errorf("\nWant: <nil>\n Got: %s", err)
	}

	if !isContract {
		t.Errorf("\nWant: true\n Got: %t", isContract)
	}
}

func TestClient_ChainID(t *testing.T) {
	tr := test.TransportRecorder(false, false, "chain_id_200_ok")

	c := newTestClient(t, "ALCHEMY", tr)
	defer c.Close()

	chainID, err := c.ChainID(context.Background())
	if err != nil {
		t.Errorf("\nWant: <nil>\n Got: %s", err)
	}

	wantChainID := "1"
	if wantChainID != chainID.String() {
		t.Errorf("\nWant: %s\n Got: %s", wantChainID, chainID.String())
	}
}

func newTestClient(t *testing.T, gatewayEnvKey string, rt http.RoundTripper) ethereum.Client { //nolint: unparam
	t.Helper()

	url, ok := os.LookupEnv(gatewayEnvKey)
	if !ok {
		t.Fatalf("gateway env key: %s not set", gatewayEnvKey)
	}

	client := &http.Client{
		Transport: rt,
		Timeout:   4 * time.Second,
	}

	c, err := ethereum.Dial(context.Background(), url, client)
	if err != nil {
		t.Fatalf("dial error: %s", err)
	}

	return c
}
