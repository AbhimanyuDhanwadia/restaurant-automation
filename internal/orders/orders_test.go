package orders

import (
	"context"
	"testing"
)

func TestServiceCreatesAndUpdatesOrder(t *testing.T) {
	service := NewService(NewMemoryRepository())
	total := int64(12550)
	order, err := service.Create(context.Background(), CreateInput{Channel: "dine-in", TotalMinor: &total, Currency: "inr", Items: []Item{{Name: "Paneer", Quantity: 1}}})
	if err != nil || order.Status != "received" {
		t.Fatalf("order=%+v err=%v", order, err)
	}
	if order.TotalMinor == nil || *order.TotalMinor != total || order.Currency != "INR" {
		t.Fatalf("sales data = %+v", order)
	}
	updated, err := service.UpdateStatus(context.Background(), order.ID, "preparing")
	if err != nil || updated.Status != "preparing" {
		t.Fatal(err)
	}
}

func TestServiceRejectsInvalidOrderTotal(t *testing.T) {
	service := NewService(NewMemoryRepository())
	total := int64(-1)
	if _, err := service.Create(context.Background(), CreateInput{Channel: "dine-in", TotalMinor: &total, Currency: "INR", Items: []Item{{Name: "Paneer", Quantity: 1}}}); err != ErrInvalidOrder {
		t.Fatalf("error = %v", err)
	}
}

func TestServiceRejectsCurrencyWithoutTotal(t *testing.T) {
	service := NewService(NewMemoryRepository())
	if _, err := service.Create(context.Background(), CreateInput{Channel: "dine-in", Currency: "INR", Items: []Item{{Name: "Paneer", Quantity: 1}}}); err != ErrInvalidOrder {
		t.Fatalf("error = %v", err)
	}
}
