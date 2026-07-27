package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/orders"
	"net/http"
)

func ListOrders(s *orders.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o, e := s.List(r.Context())
		if e != nil {
			http.Error(w, "orders unavailable", 500)
			return
		}
		writeJSON(w, 200, o)
	}
}
func CreateOrder(s *orders.Service, e *automation.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i orders.CreateInput
		if json.NewDecoder(r.Body).Decode(&i) != nil {
			http.Error(w, "invalid order", 400)
			return
		}
		o, err := s.Create(r.Context(), i)
		if err != nil {
			http.Error(w, "invalid order", 400)
			return
		}
		if e.SubmitOrder(r.Context(), o.ID) != nil {
			http.Error(w, "automation unavailable", 503)
			return
		}
		writeJSON(w, 201, o)
	}
}
func UpdateOrderStatus(s *orders.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			Status string `json:"status"`
		}
		json.NewDecoder(r.Body).Decode(&b)
		o, e := s.UpdateStatus(r.Context(), chi.URLParam(r, "orderID"), b.Status)
		if errors.Is(e, orders.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if e != nil {
			http.Error(w, "invalid status", 400)
			return
		}
		writeJSON(w, 200, o)
	}
}
