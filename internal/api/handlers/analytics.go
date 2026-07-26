package handlers

import (
	"net/http"

	"github.com/restaurantautomation/api/internal/analytics"
)

func AnalyticsOverview(service *analytics.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, service.Overview()) }
}
