package handlers

import (
	"context"
	"net/http"

	"github.com/restaurantautomation/api/internal/database"
)

type DatabaseInspector interface {
	Migrations(context.Context) ([]database.Migration, error)
}

type DatabaseStatusResponse struct {
	Status     string               `json:"status"`
	Migrations []database.Migration `json:"migrations"`
}

func DatabaseStatus(inspector DatabaseInspector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if inspector == nil {
			writeJSON(w, http.StatusOK, DatabaseStatusResponse{Status: "not_configured", Migrations: []database.Migration{}})
			return
		}
		migrations, err := inspector.Migrations(r.Context())
		if err != nil {
			http.Error(w, "database metadata unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, DatabaseStatusResponse{Status: "available", Migrations: migrations})
	}
}
