package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/printqueue"
)

func Printers(manager *printers.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, manager.Health()) }
}
func PrintQueue(queue *printqueue.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if queue == nil {
			http.Error(w, "print queue unavailable", http.StatusServiceUnavailable)
			return
		}
		jobs, err := queue.List(r.Context())
		if err != nil {
			http.Error(w, "print queue unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, jobs)
	}
}
func PrintTicket(manager *printers.Manager, queue *printqueue.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ticket printers.Ticket
		if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil || ticket.OrderID == "" || ticket.Destination == "" || len(ticket.Lines) == 0 {
			http.Error(w, "order_id, destination, and lines are required", http.StatusBadRequest)
			return
		}
		if queue == nil {
			http.Error(w, "print queue unavailable", http.StatusServiceUnavailable)
			return
		}
		job, err := queue.Queue(r.Context(), ticket)
		if err != nil {
			http.Error(w, "create print job", http.StatusInternalServerError)
			return
		}
		ticket.PrintJobID = job.ID
		if err := manager.Print(r.Context(), ticket); err != nil {
			queue.Failed(r.Context(), ticket, 0, err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusAccepted, job)
	}
}
func ReprintTicket(manager *printers.Manager, queue *printqueue.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "orderID")
		if queue == nil {
			http.Error(w, "print queue unavailable", http.StatusServiceUnavailable)
			return
		}
		job, err := queue.Reprint(r.Context(), orderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		ticket := printers.Ticket{OrderID: job.OrderID, Destination: job.Destination, Lines: job.Lines, Reprint: true, PrintJobID: job.ID}
		if err := manager.Print(r.Context(), ticket); err != nil {
			queue.Failed(r.Context(), ticket, 0, err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusAccepted, job)
	}
}
func RetryPrintJob(manager *printers.Manager, queue *printqueue.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobID")
		if queue == nil {
			http.Error(w, "print queue unavailable", http.StatusServiceUnavailable)
			return
		}
		job, err := queue.RetryFailed(r.Context(), jobID)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, printqueue.ErrNotFound) {
				status = http.StatusNotFound
			} else if errors.Is(err, printqueue.ErrJobNotFailed) {
				status = http.StatusConflict
			}
			http.Error(w, err.Error(), status)
			return
		}
		ticket := printers.Ticket{OrderID: job.OrderID, Destination: job.Destination, Lines: job.Lines, Reprint: true, PrintJobID: job.ID}
		if err := manager.Print(r.Context(), ticket); err != nil {
			queue.Failed(r.Context(), ticket, 0, err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusAccepted, job)
	}
}
