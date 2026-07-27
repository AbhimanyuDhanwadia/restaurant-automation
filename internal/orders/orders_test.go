package orders

import (
	"context"
	"testing"
)

func TestServiceCreatesAndUpdatesOrder(t *testing.T) {
	service := NewService(NewMemoryRepository())
	order, err := service.Create(context.Background(), CreateInput{Channel: "dine-in", Items: []Item{{Name: "Paneer", Quantity: 1}}})
	if err != nil || order.Status != "received" {
		t.Fatalf("order=%+v err=%v", order, err)
	}
	updated, err := service.UpdateStatus(context.Background(), order.ID, "preparing")
	if err != nil || updated.Status != "preparing" {
		t.Fatal(err)
	}
}
