package integrations

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestWebhookProviderValidatesAndEmitsOrder(t *testing.T) {
	provider := NewWebhookProvider("partner", "secret", 1)
	body := []byte(`{"order_id":"ORD-600","payload":{"channel":"delivery"}}`)
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write(body)
	if err := provider.ReceiveWebhook(body, hex.EncodeToString(mac.Sum(nil))); err != nil {
		t.Fatal(err)
	}
	order := <-provider.ReceiveOrders()
	if order.ID != "ORD-600" || order.Provider != "partner" {
		t.Fatalf("order = %+v", order)
	}
}

func TestWebhookProviderRejectsInvalidSignature(t *testing.T) {
	provider := NewWebhookProvider("partner", "secret", 1)
	if err := provider.ReceiveWebhook([]byte(`{"order_id":"ORD-600"}`), "bad"); err != ErrInvalidWebhookSignature {
		t.Fatalf("error = %v", err)
	}
}
