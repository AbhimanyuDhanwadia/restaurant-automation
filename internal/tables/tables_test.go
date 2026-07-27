package tables

import (
	"context"
	"errors"
	"testing"
)

func TestServiceCreateValidation(t *testing.T) {
	service := NewService(NewMemoryRepository())
	tests := []struct {
		name  string
		input CreateInput
		want  error
	}{
		{name: "valid", input: CreateInput{ID: "T1", Seats: 4}, want: nil},
		{name: "missing id", input: CreateInput{Seats: 4}, want: ErrInvalidTable},
		{name: "invalid capacity", input: CreateInput{ID: "T1"}, want: ErrInvalidTable},
		{name: "invalid status", input: CreateInput{ID: "T1", Seats: 4, Status: "closed"}, want: ErrInvalidStatus},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table, err := service.Create(context.Background(), test.input)
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
			if test.want == nil && table.Status != "available" {
				t.Fatalf("status = %q, want available", table.Status)
			}
		})
	}
}

func TestServiceUpdatesTableStatus(t *testing.T) {
	service := NewService(NewMemoryRepository())
	table, err := service.Create(context.Background(), CreateInput{ID: "T2", Seats: 2})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := service.UpdateStatus(context.Background(), table.ID, "seated")
	if err != nil || updated.Status != "seated" {
		t.Fatalf("table = %+v, error = %v", updated, err)
	}
}
