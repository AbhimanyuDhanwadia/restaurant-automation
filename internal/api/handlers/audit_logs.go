package handlers

import (
	"context"
	"net/http"

	"github.com/restaurantautomation/api/internal/automation"
)

type OperationalEventReader interface {
	List(context.Context, int) ([]automation.Event, error)
}

type AuditLogsResponse struct {
	Source string             `json:"source"`
	Events []automation.Event `json:"events"`
}

func AuditLogs(engine *automation.Engine, reader OperationalEventReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if reader == nil {
			writeJSON(w, http.StatusOK, AuditLogsResponse{Source: "runtime", Events: engine.Events()})
			return
		}
		events, err := reader.List(r.Context(), 200)
		if err != nil {
			http.Error(w, "audit logs unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, AuditLogsResponse{Source: "persistent", Events: events})
	}
}
