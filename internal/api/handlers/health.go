package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

type ReadinessChecker interface{ Ping(context.Context) error }

// HealthResponse represents the standard JSON response for the health endpoints.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health is a basic liveness probe. It returns 200 OK as long as the HTTP
// server is accepting connections.
func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	}
}

// Ready is a readiness probe. In Milestone 4, this will check database
// connectivity. For now, it behaves identical to Health.
func Ready() http.HandlerFunc {
	return ReadyWithCheck(nil)
}

func ReadyWithCheck(checker ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checker != nil && checker.Ping(r.Context()) != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(HealthResponse{Status: "not ready"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{Status: "ready"})
	}
}
