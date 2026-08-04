package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	appm "github.com/restaurantautomation/api/internal/api/middleware"
	"github.com/restaurantautomation/api/internal/config"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "manager-subject", "email": "manager@example.com", "exp": time.Now().Add(time.Minute).Unix()})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	updateResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPatch, "/user-1/role", bytes.NewBufferString(`{"role_id":"00000000-0000-0000-0000-000000000002"}`))
	request.Header.Set("Authorization", "Bearer "+signed)
	appm.Auth(config.AuthConfig{Required: true, SupabaseJWTSecret: "secret"})(router).ServeHTTP(updateResponse, request)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", updateResponse.Code, updateResponse.Body.String())
	}
	if events, err := service.ListRoleAssignments(context.Background(), 10); err != nil || len(events) != 1 || events[0].ActorSubject != "manager-subject" {
		t.Fatalf("events = %#v error = %v", events, err)
	}
}
