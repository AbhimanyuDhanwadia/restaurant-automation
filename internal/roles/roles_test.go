package roles

import (
	"context"
	"errors"
	"testing"
)

func TestServiceCreatesRoleWithNormalizedPermissions(t *testing.T) {
	service := NewService(NewMemoryRepository())
	role, err := service.Create(context.Background(), CreateInput{Name: "Night manager", Description: "Closes the restaurant.", Permissions: []string{"orders.manage", "analytics.view", "orders.manage"}})
	if err != nil {
		t.Fatal(err)
	}
	if role.System || len(role.Permissions) != 2 || role.Permissions[0] != "analytics.view" {
		t.Fatalf("role = %+v", role)
	}
}

func TestServiceRejectsInvalidOrDuplicateRole(t *testing.T) {
	service := NewService(NewMemoryRepository())
	if _, err := service.Create(context.Background(), CreateInput{Name: "", Permissions: []string{"orders.manage"}}); !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("invalid role error = %v", err)
	}
	if _, err := service.Create(context.Background(), CreateInput{Name: "Operator", Permissions: []string{"orders.manage"}}); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate role error = %v", err)
	}
}
