package handlers

import (
	"net/http"

	appm "github.com/restaurantautomation/api/internal/api/middleware"
	"github.com/restaurantautomation/api/internal/users"
)

type currentUserResponse struct {
	Subject     string   `json:"subject"`
	Email       string   `json:"email"`
	AuthRole    string   `json:"auth_role"`
	RoleID      string   `json:"role_id"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions"`
}

func CurrentUser(service *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := appm.ClaimsFromContext(r.Context())
		if !ok {
			http.Error(w, "authentication is not enabled", http.StatusNotImplemented)
			return
		}
		if service == nil {
			http.Error(w, "authorization unavailable", http.StatusServiceUnavailable)
			return
		}
		access, err := service.Access(r.Context(), claims.Subject)
		if err != nil {
			http.Error(w, "authorization unavailable", http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusOK, currentUserResponse{Subject: claims.Subject, Email: claims.Email, AuthRole: claims.Role, RoleID: access.User.RoleID, RoleName: access.User.RoleName, Permissions: access.Permissions})
	}
}
