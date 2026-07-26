package integrations

import (
	"context"
	"testing"
	"time"
)

func TestRegistryCollectsFromProvider(t *testing.T) {
	registry := NewRegistry()
	provider := NewMockProvider("mock", 1)
	if err := registry.Register(provider); err != nil {
		t.Fatal(err)
	}
	received := make(chan Order, 1)
	if err := registry.StartCollectors(context.Background(), func(_ context.Context, order Order) error { received <- order; return nil }); err != nil {
		t.Fatal(err)
	}
	if !provider.Submit(Order{ID: "ORD-200", Provider: provider.Name()}) {
		t.Fatal("mock provider rejected order")
	}
	select {
	case order := <-received:
		if order.ID != "ORD-200" {
			t.Fatalf("order id = %s", order.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for collected order")
	}
	provider.Close()
}

func TestRegistryRejectsDuplicateProvider(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(NewMockProvider("mock", 1)); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(NewMockProvider("mock", 1)); err != ErrProviderExists {
		t.Fatalf("error = %v", err)
	}
}
