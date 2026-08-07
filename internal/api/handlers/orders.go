package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/printqueue"
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
func CreateOrder(s *orders.Service, e *automation.Engine, manager *printers.Manager, queue *printqueue.Service) http.HandlerFunc {
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
		if err := queueKitchenTicket(r, manager, queue, o); err != nil {
			http.Error(w, "create kitchen print job", http.StatusServiceUnavailable)
			return
		}
		if e.SubmitOrder(r.Context(), o.ID) != nil {
			http.Error(w, "automation unavailable", 503)
			return
		}
		writeJSON(w, 201, o)
	}
}

// queueKitchenTicket records the ticket before handing it to the asynchronous
// printer manager. A dispatch failure is retained as a failed job so operators
// can reprint it without recreating the order.
func queueKitchenTicket(r *http.Request, manager *printers.Manager, queue *printqueue.Service, order orders.Order) error {
	if manager == nil || queue == nil {
		return nil
	}

	lines := make([]printers.Line, 0, len(order.Items))
	for _, item := range order.Items {
		lines = append(lines, printers.Line{Text: item.Name, Quantity: item.Quantity})
	}
	ticket := printers.Ticket{OrderID: order.ID, Destination: "kitchen", Lines: lines}
	job, err := queue.Queue(r.Context(), ticket)
	if err != nil {
		return err
	}
	ticket.PrintJobID = job.ID
	if err := manager.Print(r.Context(), ticket); err != nil {
		queue.Failed(r.Context(), ticket, 0, err)
	}
	return nil
}
func UpdateOrderStatus(s *orders.Service, engine *automation.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			Status string `json:"status"`
		}
		json.NewDecoder(r.Body).Decode(&b)
		o, changed, e := s.TransitionStatus(r.Context(), chi.URLParam(r, "orderID"), b.Status)
		if errors.Is(e, orders.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if e != nil {
			http.Error(w, "invalid status", 400)
			return
		}
		if changed && engine != nil {
			// The order update is authoritative. A stopped in-memory engine must
			// not cause a client to repeat an already-completed status update.
			_ = engine.RecordOrderStatus(r.Context(), o.ID, o.Status)
		}
		writeJSON(w, 200, o)
	}
}
