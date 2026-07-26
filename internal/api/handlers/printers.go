package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/printers"
)

func Printers(manager *printers.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, manager.Health()) }
}
func PrintTicket(manager *printers.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ticket printers.Ticket
		if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil || ticket.OrderID == "" || ticket.Destination == "" {
			http.Error(w, "order_id, destination, and lines are required", http.StatusBadRequest)
			return
		}
		if err := manager.Print(r.Context(), ticket); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "queued", "order_id": ticket.OrderID})
	}
}
func ReprintTicket(manager *printers.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "orderID")
		if err := manager.Reprint(r.Context(), orderID); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "reprint_queued", "order_id": orderID})
	}
}
