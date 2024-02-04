package main

import (
	"fmt"
	"os"
	"time"
)

type config struct {
	Address                string
	Gateways               map[string]string
	GatewayTimeout         time.Duration
	GatewaysMinimumClients int
}

func newConfig(address string) (config, error) {
	gateways, err := lookupEnvGateways()
	if err != nil {
		return config{}, err
	}

	cfg := config{
		Address:                address,
		Gateways:               gateways,
		GatewayTimeout:         5 * time.Second,
		GatewaysMinimumClients: 4,
	}

	return cfg, nil
}

func lookupEnvGateways() (map[string]string, error) {
	urlMap := map[string]string{}

	if url, ok := os.LookupEnv("ALCHEMY"); ok {
		urlMap["Alchemy"] = url
	}

	if url, ok := os.LookupEnv("INFURA"); ok {
		urlMap["Infura"] = url
	}

	if url, ok := os.LookupEnv("CHAINSTACK"); ok {
		urlMap["Chainstack"] = url
	}

	if url, ok := os.LookupEnv("TENDERLY"); ok {
		urlMap["Tenderly"] = url
	}

	const networks = 4
	if len(urlMap) != networks {
		return urlMap, fmt.Errorf("you must configure four Ethereum networks in the .env file")
	}

	return urlMap, nil
}
