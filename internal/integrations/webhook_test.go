package integrations

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

func TestWebhookProviderValidatesAndEmitsOrder(t *testing.T) {
	provider := NewWebhookProvider("partner", "secret", 1)
	body := []byte(`{"order_id":"ORD-600","payload":{"channel":"delivery"}}`)
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write(body)
	if err := provider.ReceiveWebhook(context.Background(), body, hex.EncodeToString(mac.Sum(nil))); err != nil {
		t.Fatal(err)
	}
	order := <-provider.ReceiveOrders()
	if order.ID != "ORD-600" || order.Provider != "partner" {
		t.Fatalf("order = %+v", order)
	}
}

func TestWebhookProviderRejectsInvalidSignature(t *testing.T) {
	provider := NewWebhookProvider("partner", "secret", 1)
	if err := provider.ReceiveWebhook(context.Background(), []byte(`{"order_id":"ORD-600"}`), "bad"); err != ErrInvalidWebhookSignature {
		t.Fatalf("error = %v", err)
	}
}

func TestWebhookProviderDeduplicatesAcceptedOrder(t *testing.T) {
	provider := NewWebhookProvider("partner", "secret", 2)
	body := []byte(`{"order_id":"ORD-601"}`)
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))
	if err := provider.ReceiveWebhook(context.Background(), body, signature); err != nil {
		t.Fatal(err)
	}
	if err := provider.ReceiveWebhook(context.Background(), body, signature); err != ErrWebhookDuplicate {
		t.Fatalf("error = %v, want duplicate", err)
	}
	if order := <-provider.ReceiveOrders(); order.ID != "ORD-601" {
		t.Fatalf("order = %#v", order)
	}
	select {
	case order := <-provider.ReceiveOrders():
		t.Fatalf("unexpected duplicate order = %#v", order)
	default:
	}
}

func TestWebhookProviderSharesReceiptStoreAcrossInstances(t *testing.T) {
	receipts := NewMemoryWebhookReceiptStore()
	first := NewWebhookProviderWithReceipts("partner", "secret", 1, receipts)
	second := NewWebhookProviderWithReceipts("partner", "secret", 1, receipts)
	body := []byte(`{"order_id":"ORD-602"}`)
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))
	if err := first.ReceiveWebhook(context.Background(), body, signature); err != nil {
		t.Fatal(err)
	}
	if err := second.ReceiveWebhook(context.Background(), body, signature); err != ErrWebhookDuplicate {
		t.Fatalf("error = %v, want duplicate", err)
	}
}

func TestWebhookProviderReleasesReceiptWhenQueueIsFull(t *testing.T) {
	receipts := NewMemoryWebhookReceiptStore()
	provider := NewWebhookProviderWithReceipts("partner", "secret", 1, receipts)
	sign := func(orderID string) ([]byte, string) {
		body := []byte(`{"order_id":"` + orderID + `"}`)
		mac := hmac.New(sha256.New, []byte("secret"))
		mac.Write(body)
		return body, hex.EncodeToString(mac.Sum(nil))
	}
	firstBody, firstSignature := sign("ORD-603")
	if err := provider.ReceiveWebhook(context.Background(), firstBody, firstSignature); err != nil {
		t.Fatal(err)
	}
	secondBody, secondSignature := sign("ORD-604")
	if err := provider.ReceiveWebhook(context.Background(), secondBody, secondSignature); err != ErrWebhookQueueFull {
		t.Fatalf("error = %v, want queue full", err)
	}
	retryProvider := NewWebhookProviderWithReceipts("partner", "secret", 1, receipts)
	if err := retryProvider.ReceiveWebhook(context.Background(), secondBody, secondSignature); err != nil {
		t.Fatalf("receipt was not released: %v", err)
	}
}

func TestWebhookProviderMarksConcurrentReceiptAsInProgress(t *testing.T) {
	receipts := NewMemoryWebhookReceiptStore()
	now := time.Now().UTC()
	if claim, err := receipts.Claim(context.Background(), "partner", "ORD-605", now, now.Add(webhookDedupeWindow)); err != nil || claim != WebhookReceiptClaimed {
		t.Fatalf("claim = %q error = %v", claim, err)
	}
	provider := NewWebhookProviderWithReceipts("partner", "secret", 1, receipts)
	body := []byte(`{"order_id":"ORD-605"}`)
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write(body)
	if err := provider.ReceiveWebhook(context.Background(), body, hex.EncodeToString(mac.Sum(nil))); err != ErrWebhookInProgress {
		t.Fatalf("error = %v, want in progress", err)
	}
}
