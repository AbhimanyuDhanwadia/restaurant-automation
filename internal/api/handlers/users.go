package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/users"
)

type updateUserRoleRequest struct {
	RoleID string `json:"role_id"`
}

func ListUsers(service *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "users unavailable", http.StatusServiceUnavailable)
			return
		}
		result, err := service.List(r.Context())
		if err != nil {
			http.Error(w, "users unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func UpdateUserRole(service *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "users unavailable", http.StatusServiceUnavailable)
			return
		}
		var request updateUserRoleRequest
		if json.NewDecoder(r.Body).Decode(&request) != nil {
			http.Error(w, "invalid user role", http.StatusBadRequest)
			return
		}
		user, err := service.UpdateRole(r.Context(), chi.URLParam(r, "subject"), request.RoleID)
		if errors.Is(err, users.ErrNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, users.ErrInvalidUser) || errors.Is(err, users.ErrInvalidRole) {
			http.Error(w, "invalid user role", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "users unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}
