package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/restaurantautomation/api/internal/backups"
)

func ListBackups(service *backups.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "backups unavailable", http.StatusServiceUnavailable)
			return
		}
		records, err := service.List(r.Context())
		if err != nil {
			http.Error(w, "backups unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, records)
	}
}

func RecordBackup(service *backups.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "backups unavailable", http.StatusServiceUnavailable)
			return
		}
		var input backups.RecordInput
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			http.Error(w, "invalid backup record", http.StatusBadRequest)
			return
		}
		record, err := service.Record(r.Context(), input)
		if err != nil {
			http.Error(w, "invalid backup record", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, record)
	}
}
