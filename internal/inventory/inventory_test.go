package inventory

import (
	"context"
	"errors"
	"testing"
)

func TestServiceStockStatus(t *testing.T) {
	service := NewService(NewMemoryRepository())
	item, err := service.Create(context.Background(), CreateInput{Name: "Tomatoes", Category: "Produce", OnHand: 2, ParLevel: 4, Unit: "kg"})
	if err != nil || item.Status != "low_stock" {
		t.Fatalf("item = %+v, err = %v", item, err)
	}
	updated, err := service.UpdateStock(context.Background(), item.ID, 6)
	if err != nil || updated.Status != "in_stock" {
		t.Fatalf("item = %+v, err = %v", updated, err)
	}
}

func TestServiceValidatesItem(t *testing.T) {
	service := NewService(NewMemoryRepository())
	_, err := service.Create(context.Background(), CreateInput{Name: "Tomatoes", Category: "Produce", OnHand: -1, Unit: "kg"})
	if !errors.Is(err, ErrInvalidItem) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidItem)
	}
}
