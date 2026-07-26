package integrations

import (
	"context"
	"errors"
	"sort"
	"sync"
)

type Status string

const (
	StatusConnected    Status = "connected"
	StatusDisconnected Status = "disconnected"
	StatusError        Status = "error"
)

type Order struct {
	ID       string         `json:"id"`
	Provider string         `json:"provider"`
	Payload  map[string]any `json:"payload,omitempty"`
}

// Provider is the only contract the automation engine needs from an order source.
type Provider interface {
	Name() string
	Connect() error
	Disconnect() error
	Health() Status
	ReceiveOrders() <-chan Order
}

var ErrProviderExists = errors.New("provider is already registered")
var ErrProviderNotFound = errors.New("provider is not registered")

type Health struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}

type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

func NewRegistry() *Registry { return &Registry{providers: make(map[string]Provider)} }

func (r *Registry) Register(provider Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.providers[provider.Name()]; exists {
		return ErrProviderExists
	}
	r.providers[provider.Name()] = provider
	return nil
}

func (r *Registry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.providers[name]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return provider, nil
}

func (r *Registry) Snapshot() []Health {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Health, 0, len(r.providers))
	for _, provider := range r.providers {
		result = append(result, Health{Name: provider.Name(), Status: provider.Health()})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func (r *Registry) ConnectAll() error {
	r.mu.RLock()
	providers := make([]Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	r.mu.RUnlock()
	for _, provider := range providers {
		if err := provider.Connect(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) DisconnectAll() error {
	r.mu.RLock()
	providers := make([]Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	r.mu.RUnlock()
	for _, provider := range providers {
		if err := provider.Disconnect(); err != nil {
			return err
		}
	}
	return nil
}

// StartCollectors forwards provider orders to the supplied sink. Provider-specific
// payloads stop at this boundary and can be normalized by the automation engine.
func (r *Registry) StartCollectors(ctx context.Context, sink func(context.Context, Order) error) error {
	if err := r.ConnectAll(); err != nil {
		return err
	}
	r.mu.RLock()
	providers := make([]Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	r.mu.RUnlock()
	for _, provider := range providers {
		orders := provider.ReceiveOrders()
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case order, ok := <-orders:
					if !ok {
						return
					}
					_ = sink(ctx, order)
				}
			}
		}()
	}
	return nil
}
