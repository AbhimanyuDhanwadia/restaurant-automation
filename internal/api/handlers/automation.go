package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/restaurantautomation/api/internal/automation"
)

type submitOrderRequest struct {
	OrderID string `json:"order_id"`
}

func SubmitAutomationOrder(engine *automation.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request submitOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil || strings.TrimSpace(request.OrderID) == "" {
			http.Error(w, "order_id is required", http.StatusBadRequest)
			return
		}
		request.OrderID = strings.TrimSpace(request.OrderID)
		if err := engine.SubmitOrder(r.Context(), request.OrderID); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted", "order_id": request.OrderID})
	}
}
func AutomationEvents(engine *automation.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, engine.Events()) }
}
func AutomationQueue(engine *automation.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, engine.Stats()) }
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
