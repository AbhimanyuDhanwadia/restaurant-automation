package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/alerts"
)

func ListAlerts(service *alerts.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := service.ListActive(r.Context())
		if err != nil {
			http.Error(w, "alerts unavailable", 500)
			return
		}
		writeJSON(w, 200, result)
	}
}
func CreateAlert(service *alerts.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input alerts.CreateInput
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			http.Error(w, "invalid alert", 400)
			return
		}
		alert, err := service.Create(r.Context(), input)
		if err != nil {
			http.Error(w, "invalid alert", 400)
			return
		}
		writeJSON(w, 201, alert)
	}
}
func AcknowledgeAlert(service *alerts.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := service.Acknowledge(r.Context(), chi.URLParam(r, "alertID"))
		if errors.Is(err, alerts.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "alert unavailable", 500)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
