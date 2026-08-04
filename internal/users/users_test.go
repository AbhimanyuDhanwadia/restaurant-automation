package users

import (
	"context"
	"errors"
	"testing"

	"github.com/restaurantautomation/api/internal/roles"
)

func newTestService(adminEmails ...string) *Service {
	return NewService(NewMemoryRepository(), roles.NewService(roles.NewMemoryRepository()), adminEmails)
}

func TestObserveAssignsBootstrapAdministratorAndPreservesAssignments(t *testing.T) {
	service := newTestService("Manager@example.com")
	if err := service.Observe(context.Background(), "manager-subject", "MANAGER@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := service.Observe(context.Background(), "operator-subject", "operator@example.com"); err != nil {
		t.Fatal(err)
	}
	listed, err := service.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 2 || listed[0].Email != "manager@example.com" || listed[0].RoleName != "Administrator" || listed[1].RoleName != "Operator" {
		t.Fatalf("users = %#v", listed)
	}
	if _, err := service.UpdateRole(context.Background(), "manager-subject", roles.ManagerID); err != nil {
		t.Fatal(err)
	}
	if err := service.Observe(context.Background(), "manager-subject", "manager@example.com"); err != nil {
		t.Fatal(err)
	}
	listed, err = service.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if listed[0].RoleName != "Manager" {
		t.Fatalf("role = %q, want Manager", listed[0].RoleName)
	}
}

func TestUpdateRoleRejectsUnknownRoleAndUnknownUser(t *testing.T) {
	service := newTestService()
	if err := service.Observe(context.Background(), "operator-subject", "operator@example.com"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.UpdateRole(context.Background(), "operator-subject", "not-a-role"); !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("error = %v, want invalid role", err)
	}
	if _, err := service.UpdateRole(context.Background(), "missing", roles.ManagerID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("error = %v, want not found", err)
	}
}

func TestObserveRejectsIncompleteIdentity(t *testing.T) {
	service := newTestService()
	if err := service.Observe(context.Background(), "", "operator@example.com"); !errors.Is(err, ErrInvalidUser) {
		t.Fatalf("error = %v, want invalid user", err)
	}
}

func TestPermissionsFollowAssignedRole(t *testing.T) {
	service := newTestService()
	if err := service.Observe(context.Background(), "operator-subject", "operator@example.com"); err != nil {
		t.Fatal(err)
	}
	permissions, err := service.Permissions(context.Background(), "operator-subject")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(permissions, "orders.manage") || contains(permissions, "users.manage") {
		t.Fatalf("operator permissions = %#v", permissions)
	}
	if _, err := service.UpdateRole(context.Background(), "operator-subject", roles.AdministratorID); err != nil {
		t.Fatal(err)
	}
	permissions, err = service.Permissions(context.Background(), "operator-subject")
	if err != nil || !contains(permissions, "users.manage") || !contains(permissions, "backups.manage") {
		t.Fatalf("administrator permissions = %#v error = %v", permissions, err)
	}
	access, err := service.Access(context.Background(), "operator-subject")
	if err != nil || access.User.RoleName != "Administrator" {
		t.Fatalf("access = %#v error = %v", access, err)
	}
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
