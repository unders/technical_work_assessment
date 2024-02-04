// nolint: gochecknoglobals
package metric

import "github.com/prometheus/client_golang/prometheus"

var EthBalanceCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "eth_balance_request_count_total",
		Help: "No of request handled by handler_eth_balance",
	},
)
