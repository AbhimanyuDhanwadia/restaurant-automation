package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/integrations"
)

func ReceiveWebhook(registry *integrations.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if err != nil {
			http.Error(w, "invalid webhook body", http.StatusBadRequest)
			return
		}
		err = registry.ReceiveWebhook(chi.URLParam(r, "provider"), body, r.Header.Get("X-Webhook-Signature"))
		switch {
		case err == nil:
			writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
		case errors.Is(err, integrations.ErrWebhookDuplicate):
			writeJSON(w, http.StatusAccepted, map[string]string{"status": "duplicate"})
		case errors.Is(err, integrations.ErrProviderNotFound), errors.Is(err, integrations.ErrWebhookUnsupported):
			http.NotFound(w, r)
		case errors.Is(err, integrations.ErrInvalidWebhookSignature):
			http.Error(w, "invalid webhook signature", http.StatusUnauthorized)
		case errors.Is(err, integrations.ErrInvalidWebhookPayload):
			http.Error(w, "invalid webhook payload", http.StatusBadRequest)
		default:
			http.Error(w, "webhook unavailable", http.StatusServiceUnavailable)
		}
	}
}
