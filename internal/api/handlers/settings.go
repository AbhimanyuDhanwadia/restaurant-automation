package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/restaurantautomation/api/internal/settings"
)

func GetSettings(service *settings.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := service.Get(r.Context())
		if err != nil {
			http.Error(w, "settings unavailable", 500)
			return
		}
		writeJSON(w, 200, value)
	}
}
func SaveSettings(service *settings.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var value settings.Settings
		if json.NewDecoder(r.Body).Decode(&value) != nil {
			http.Error(w, "invalid settings", 400)
			return
		}
		saved, err := service.Save(r.Context(), value)
		if err != nil {
			http.Error(w, "invalid settings", 400)
			return
		}
		writeJSON(w, 200, saved)
	}
}
