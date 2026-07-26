package handlers

import (
	"net/http"

	"github.com/restaurantautomation/api/internal/integrations"
)

func Integrations(registry *integrations.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, registry.Snapshot()) }
}
