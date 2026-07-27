package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/tables"
)

func ListTables(service *tables.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := service.List(r.Context())
		if err != nil {
			http.Error(w, "tables unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func CreateTable(service *tables.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input tables.CreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid table", http.StatusBadRequest)
			return
		}
		table, err := service.Create(r.Context(), input)
		if err != nil {
			http.Error(w, "invalid table", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, table)
	}
}

func UpdateTableStatus(service *tables.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		table, err := service.UpdateStatus(r.Context(), chi.URLParam(r, "tableID"), body.Status)
		if errors.Is(err, tables.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, table)
	}
}
