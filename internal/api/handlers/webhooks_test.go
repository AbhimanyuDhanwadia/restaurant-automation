package handlers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/integrations"
)

func TestReceiveWebhookAcceptsSignedOrder(t *testing.T) {
	registry := integrations.NewRegistry()
	if err := registry.Register(integrations.NewWebhookProvider("partner", "secret", 1)); err != nil {
		t.Fatal(err)
	}
	body := []byte(`{"order_id":"ORD-701"}`)
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write(body)
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	request.Header.Set("X-Webhook-Signature", hex.EncodeToString(mac.Sum(nil)))
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("provider", "partner")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
	response := httptest.NewRecorder()
	ReceiveWebhook(registry).ServeHTTP(response, request)
	if response.Code != http.StatusAccepted {
		t.Fatalf("status = %d", response.Code)
	}
}
