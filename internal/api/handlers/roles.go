package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/restaurantautomation/api/internal/roles"
)

func ListRoles(service *roles.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "roles unavailable", http.StatusServiceUnavailable)
			return
		}
		result, err := service.List(r.Context())
		if err != nil {
			http.Error(w, "roles unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func CreateRole(service *roles.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "roles unavailable", http.StatusServiceUnavailable)
			return
		}
		var input roles.CreateInput
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			http.Error(w, "invalid role", http.StatusBadRequest)
			return
		}
		role, err := service.Create(r.Context(), input)
		if errors.Is(err, roles.ErrDuplicate) {
			http.Error(w, "role already exists", http.StatusConflict)
			return
		}
		if err != nil {
			http.Error(w, "invalid role", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, role)
	}
}
