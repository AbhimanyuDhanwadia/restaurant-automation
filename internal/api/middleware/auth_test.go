package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/restaurantautomation/api/internal/config"
	"github.com/restaurantautomation/api/internal/roles"
	"github.com/restaurantautomation/api/internal/users"
)

func TestAuthRejectsMissingToken(t *testing.T) {
	handler := Auth(config.AuthConfig{Required: true, SupabaseJWTSecret: "secret"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestTrackUserRecordsValidatedIdentity(t *testing.T) {
	service := users.NewService(users.NewMemoryRepository(), roles.NewService(roles.NewMemoryRepository()), nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user-1", "email": "operator@example.com", "exp": time.Now().Add(time.Minute).Unix()})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	handler := Auth(config.AuthConfig{Required: true, SupabaseJWTSecret: "secret"})(TrackUser(service)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+signed)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d", response.Code)
	}
	listed, err := service.List(context.Background())
	if err != nil || len(listed) != 1 || listed[0].Email != "operator@example.com" {
		t.Fatalf("users = %#v error = %v", listed, err)
	}
}

func TestRequirePermissionUsesAssignedRole(t *testing.T) {
	service := users.NewService(users.NewMemoryRepository(), roles.NewService(roles.NewMemoryRepository()), []string{"manager@example.com"})
	operatorToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "operator", "email": "operator@example.com", "exp": time.Now().Add(time.Minute).Unix()})
	operatorSigned, err := operatorToken.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	handler := Auth(config.AuthConfig{Required: true, SupabaseJWTSecret: "secret"})(TrackUser(service)(RequirePermission(service, "users.manage")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+operatorSigned)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d", response.Code)
	}

	adminToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "manager", "email": "manager@example.com", "exp": time.Now().Add(time.Minute).Unix()})
	adminSigned, err := adminToken.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	request = httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+adminSigned)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("admin status = %d", response.Code)
	}
}

func TestAuthAcceptsValidToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user-1", "email": "manager@example.com", "role": "authenticated", "exp": time.Now().Add(time.Minute).Unix()})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	handler := Auth(config.AuthConfig{Required: true, SupabaseJWTSecret: "secret"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok || claims.Subject != "user-1" {
			t.Fatal("claims missing from context")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+signed)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d", response.Code)
	}
}
