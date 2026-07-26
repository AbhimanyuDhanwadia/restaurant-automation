package handlers

import (
	"net/http"

	appm "github.com/restaurantautomation/api/internal/api/middleware"
)

func CurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := appm.ClaimsFromContext(r.Context())
		if !ok {
			http.Error(w, "authentication is not enabled", http.StatusNotImplemented)
			return
		}
		writeJSON(w, http.StatusOK, claims)
	}
}
