package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	appm "github.com/restaurantautomation/api/internal/api/middleware"
	"github.com/restaurantautomation/api/internal/config"
	"github.com/restaurantautomation/api/internal/roles"
	"github.com/restaurantautomation/api/internal/users"
)

func TestCurrentUserReturnsLocalAccessProfile(t *testing.T) {
	service := users.NewService(users.NewMemoryRepository(), roles.NewService(roles.NewMemoryRepository()), []string{"manager@example.com"})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "manager-subject", "email": "manager@example.com", "role": "authenticated", "exp": time.Now().Add(time.Minute).Unix()})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	handler := appm.Auth(config.AuthConfig{Required: true, SupabaseJWTSecret: "secret"})(appm.TrackUser(service)(CurrentUser(service)))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	request.Header.Set("Authorization", "Bearer "+signed)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	var result currentUserResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.RoleName != "Administrator" || !containsPermission(result.Permissions, "users.manage") {
		t.Fatalf("access profile = %#v", result)
	}
}

func containsPermission(permissions []string, expected string) bool {
	for _, permission := range permissions {
		if permission == expected {
			return true
		}
	}
	return false
}
