package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/integrations"
	"github.com/restaurantautomation/api/internal/printers"
)

type ReadinessChecker interface{ Ping(context.Context) error }

// HealthResponse represents the standard JSON response for the health endpoints.
type HealthResponse struct {
	Status string `json:"status"`
}

type SystemHealthComponent struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

type SystemHealthResponse struct {
	Status     string                  `json:"status"`
	Components []SystemHealthComponent `json:"components"`
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

func SystemHealth(engine *automation.Engine, registry *integrations.Registry, printerManager *printers.Manager, checker ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		components := []SystemHealthComponent{databaseHealth(r.Context(), checker), queueHealth(engine), integrationHealth(registry), printerHealth(printerManager)}
		status := "healthy"
		for _, component := range components {
			if component.Status == "unavailable" {
				status = "degraded"
				break
			}
		}
		writeJSON(w, http.StatusOK, SystemHealthResponse{Status: status, Components: components})
	}
}

func databaseHealth(ctx context.Context, checker ReadinessChecker) SystemHealthComponent {
	if checker == nil {
		return SystemHealthComponent{Name: "database", Status: "not_configured", Detail: "PostgreSQL is not configured"}
	}
	if checker.Ping(ctx) != nil {
		return SystemHealthComponent{Name: "database", Status: "unavailable", Detail: "Readiness check failed"}
	}
	return SystemHealthComponent{Name: "database", Status: "healthy", Detail: "Readiness check passed"}
}

func queueHealth(engine *automation.Engine) SystemHealthComponent {
	stats := engine.Stats()
	return SystemHealthComponent{Name: "queue", Status: "healthy", Detail: strconv.Itoa(stats.QueueDepth) + " waiting of " + strconv.Itoa(stats.QueueCapacity)}
}

func integrationHealth(registry *integrations.Registry) SystemHealthComponent {
	providers := registry.Snapshot()
	if len(providers) == 0 {
		return SystemHealthComponent{Name: "integrations", Status: "not_configured", Detail: "No providers are registered"}
	}
	connected := 0
	for _, provider := range providers {
		if provider.Status == integrations.StatusConnected {
			connected++
		}
	}
	if connected != len(providers) {
		return SystemHealthComponent{Name: "integrations", Status: "unavailable", Detail: strconv.Itoa(connected) + " of " + strconv.Itoa(len(providers)) + " providers connected"}
	}
	return SystemHealthComponent{Name: "integrations", Status: "healthy", Detail: strconv.Itoa(connected) + " providers connected"}
}

func printerHealth(printerManager *printers.Manager) SystemHealthComponent {
	fleet := printerManager.Health()
	if len(fleet) == 0 {
		return SystemHealthComponent{Name: "printers", Status: "not_configured", Detail: "No printers are registered"}
	}
	ready := 0
	for _, printer := range fleet {
		if printer.Status == printers.StatusReady {
			ready++
		}
	}
	if ready != len(fleet) {
		return SystemHealthComponent{Name: "printers", Status: "unavailable", Detail: strconv.Itoa(ready) + " of " + strconv.Itoa(len(fleet)) + " printers ready"}
	}
	return SystemHealthComponent{Name: "printers", Status: "healthy", Detail: strconv.Itoa(ready) + " printers ready"}
}
