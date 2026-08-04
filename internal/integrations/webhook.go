package integrations

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrWebhookUnsupported = errors.New("provider does not accept webhooks")
var ErrInvalidWebhookSignature = errors.New("invalid webhook signature")
var ErrInvalidWebhookPayload = errors.New("invalid webhook payload")
var ErrWebhookQueueFull = errors.New("webhook order queue is full")
var ErrWebhookDuplicate = errors.New("duplicate webhook order")

const webhookDedupeWindow = 15 * time.Minute
const maxRememberedWebhookOrders = 10_000

type WebhookProvider struct {
	name   string
	secret []byte
	orders chan Order
	mu     sync.RWMutex
	status Status

	dedupeMu sync.Mutex
	accepted map[string]time.Time
	now      func() time.Time
}

func NewWebhookProvider(name, secret string, buffer int) *WebhookProvider {
	if buffer < 1 {
		buffer = 1
	}
	return &WebhookProvider{name: name, secret: []byte(secret), orders: make(chan Order, buffer), status: StatusDisconnected, accepted: make(map[string]time.Time), now: time.Now}
}
func (p *WebhookProvider) Name() string { return p.name }
func (p *WebhookProvider) Connect() error {
	p.mu.Lock()
	p.status = StatusConnected
	p.mu.Unlock()
	return nil
}
func (p *WebhookProvider) Disconnect() error {
	p.mu.Lock()
	p.status = StatusDisconnected
	p.mu.Unlock()
	return nil
}
func (p *WebhookProvider) Health() Status              { p.mu.RLock(); defer p.mu.RUnlock(); return p.status }
func (p *WebhookProvider) ReceiveOrders() <-chan Order { return p.orders }

func (p *WebhookProvider) ReceiveWebhook(body []byte, signature string) error {
	if !p.validSignature(body, signature) {
		return ErrInvalidWebhookSignature
	}
	var payload struct {
		OrderID string         `json:"order_id"`
		Payload map[string]any `json:"payload"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || strings.TrimSpace(payload.OrderID) == "" {
		return ErrInvalidWebhookPayload
	}
	order := Order{ID: strings.TrimSpace(payload.OrderID), Provider: p.name, Payload: payload.Payload}
	return p.enqueueOrder(order)
}

func (p *WebhookProvider) enqueueOrder(order Order) error {
	now := p.now().UTC()
	p.dedupeMu.Lock()
	defer p.dedupeMu.Unlock()
	p.pruneAccepted(now)
	if _, exists := p.accepted[order.ID]; exists {
		return ErrWebhookDuplicate
	}
	select {
	case p.orders <- order:
		p.accepted[order.ID] = now
		return nil
	default:
		return ErrWebhookQueueFull
	}
}

func (p *WebhookProvider) pruneAccepted(now time.Time) {
	cutoff := now.Add(-webhookDedupeWindow)
	oldestID := ""
	var oldest time.Time
	for orderID, acceptedAt := range p.accepted {
		if !acceptedAt.After(cutoff) {
			delete(p.accepted, orderID)
			continue
		}
		if oldestID == "" || acceptedAt.Before(oldest) {
			oldestID = orderID
			oldest = acceptedAt
		}
	}
	if len(p.accepted) >= maxRememberedWebhookOrders && oldestID != "" {
		delete(p.accepted, oldestID)
	}
}

func (p *WebhookProvider) validSignature(body []byte, signature string) bool {
	signature = strings.TrimPrefix(strings.TrimSpace(signature), "sha256=")
	expected := hmac.New(sha256.New, p.secret)
	expected.Write(body)
	supplied, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(expected.Sum(nil), supplied) == 1
}

type WebhookReceiver interface{ ReceiveWebhook([]byte, string) error }

func (r *Registry) ReceiveWebhook(providerName string, body []byte, signature string) error {
	provider, err := r.Get(providerName)
	if err != nil {
		return err
	}
	receiver, ok := provider.(WebhookReceiver)
	if !ok {
		return ErrWebhookUnsupported
	}
	return receiver.ReceiveWebhook(body, signature)
}
