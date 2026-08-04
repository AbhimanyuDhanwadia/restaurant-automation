package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/roles"
	"github.com/restaurantautomation/api/internal/users"
)

func TestUserHandlersListAndUpdateRole(t *testing.T) {
	service := users.NewService(users.NewMemoryRepository(), roles.NewService(roles.NewMemoryRepository()), nil)
	if err := service.Observe(context.Background(), "user-1", "operator@example.com"); err != nil {
		t.Fatal(err)
	}
	listResponse := httptest.NewRecorder()
	ListUsers(service).ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil))
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d", listResponse.Code)
	}

	router := chi.NewRouter()
	router.Patch("/{subject}/role", UpdateUserRole(service))
	updateResponse := httptest.NewRecorder()
	router.ServeHTTP(updateResponse, httptest.NewRequest(http.MethodPatch, "/user-1/role", bytes.NewBufferString(`{"role_id":"00000000-0000-0000-0000-000000000002"}`)))
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", updateResponse.Code, updateResponse.Body.String())
	}
}
