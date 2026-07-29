package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/inventory"
)

func ListInventory(service *inventory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := service.List(r.Context())
		if err != nil {
			http.Error(w, "inventory unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

func CreateInventoryItem(service *inventory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input inventory.CreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid inventory item", http.StatusBadRequest)
			return
		}
		item, err := service.Create(r.Context(), input)
		if err != nil {
			http.Error(w, "invalid inventory item", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	}
}

func UpdateInventoryStock(service *inventory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			OnHand *float64 `json:"on_hand"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.OnHand == nil {
			http.Error(w, "invalid stock quantity", http.StatusBadRequest)
			return
		}
		item, err := service.UpdateStock(r.Context(), chi.URLParam(r, "itemID"), *body.OnHand)
		if errors.Is(err, inventory.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "invalid stock quantity", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, item)
	}
}

func UpdateInventoryStatus(service *inventory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid inventory status", http.StatusBadRequest)
			return
		}
		item, err := service.UpdateStatus(r.Context(), chi.URLParam(r, "itemID"), body.Status)
		if errors.Is(err, inventory.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "invalid inventory status", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, item)
	}
}
