package handlers

import (
	"net/http"

	"github.com/restaurantautomation/api/internal/intelligence"
)

func Insights(service *intelligence.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, service.Insights()) }
}
