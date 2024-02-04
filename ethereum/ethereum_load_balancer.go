package ethereum

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
)

type Options struct {
	MinClients int
	Gateways   map[string]string
	HTTPClient *http.Client
	Log        *slog.Logger
}

type LoadBalancer struct {
	clients []Client

	liveMtx     sync.RWMutex
	liveClients []Client

	log *slog.Logger
}

func NewLoadBalancer(ctx context.Context, opts Options) (*LoadBalancer, error) {
	clients, err := dialClients(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("dial clients error: %w", err)
	}

	lb := &LoadBalancer{clients: clients, liveMtx: sync.RWMutex{}, log: opts.Log}

	if ok := lb.HealthCheck(ctx); !ok {
		return lb, errors.New("zero ethereum backends is up")
	}

	return lb, nil
}

func (lb *LoadBalancer) HealthCheck(ctx context.Context) bool {
	liveClients := make([]Client, 0, len(lb.clients))

	for _, c := range lb.clients {
		if _, err := c.ChainID(ctx); err != nil {
			lb.log.Warn(err.Error())

			continue
		}

		liveClients = append(liveClients, c)
	}

	if len(liveClients) < 1 {
		return false
	}

	lb.liveMtx.Lock()
	lb.liveClients = liveClients
	lb.liveMtx.Unlock()

	return true
}

func (lb *LoadBalancer) Call(f func(c Client) error) error {
	lb.liveMtx.RLock()
	if len(lb.liveClients) < 1 {
		lb.liveMtx.RUnlock()

		return errors.New("zero ethereum backends is up")
	}

	liveClients := make([]Client, len(lb.liveClients))
	copy(liveClients, lb.liveClients)
	lb.liveMtx.RUnlock()

	var err error

	size := len(liveClients)
	startIndex := rand.Intn(size) // #nosec

	for i := startIndex; i < size; {
		err = f(liveClients[i])
		if err == nil {
			return nil
		}

		lb.log.Warn(err.Error())

		i++
		if i == size {
			i = 0
		}

		if i == startIndex {
			break
		}
	}

	return fmt.Errorf("all ethereum backends returned an error; last backend error: %w", err)
}

func (lb *LoadBalancer) Close() {
	for _, c := range lb.clients {
		c.Close()
	}
}

func dialClients(ctx context.Context, opts Options) ([]Client, error) {
	var (
		clients = make([]Client, 0, len(opts.Gateways))
		log     = opts.Log
	)

	for _, url := range opts.Gateways {
		ethClient, err := Dial(ctx, url, opts.HTTPClient)
		if err != nil {
			log.Error(err.Error())

			continue
		}

		clients = append(clients, ethClient)
	}

	if len(clients) < opts.MinClients {
		return clients, fmt.Errorf("less that %d ethereum clients", opts.MinClients)
	}

	return clients, nil
}
