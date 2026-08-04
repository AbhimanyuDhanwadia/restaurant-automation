package integrations

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrWebhookReceiptUnavailable = errors.New("webhook receipt store unavailable")

type WebhookReceiptClaim string

const (
	WebhookReceiptClaimed    WebhookReceiptClaim = "claimed"
	WebhookReceiptDuplicate  WebhookReceiptClaim = "duplicate"
	WebhookReceiptInProgress WebhookReceiptClaim = "in_progress"
)

// WebhookReceiptStore atomically claims a provider order ID for a bounded
// window and distinguishes accepted receipts from in-flight claims.
type WebhookReceiptStore interface {
	Claim(context.Context, string, string, time.Time, time.Time) (WebhookReceiptClaim, error)
	Confirm(context.Context, string, string) error
	Release(context.Context, string, string) error
}

type MemoryWebhookReceiptStore struct {
	mu       sync.Mutex
	receipts map[string]memoryWebhookReceipt
}

type memoryWebhookReceipt struct {
	expiresAt time.Time
	accepted  bool
}

func NewMemoryWebhookReceiptStore() *MemoryWebhookReceiptStore {
	return &MemoryWebhookReceiptStore{receipts: make(map[string]memoryWebhookReceipt)}
}

func (s *MemoryWebhookReceiptStore) Claim(_ context.Context, provider, orderID string, now, expiresAt time.Time) (WebhookReceiptClaim, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, receipt := range s.receipts {
		if !receipt.expiresAt.After(now) {
			delete(s.receipts, key)
		}
	}
	key := provider + "\x00" + orderID
	if receipt, exists := s.receipts[key]; exists {
		if receipt.accepted {
			return WebhookReceiptDuplicate, nil
		}
		return WebhookReceiptInProgress, nil
	}
	s.receipts[key] = memoryWebhookReceipt{expiresAt: expiresAt}
	return WebhookReceiptClaimed, nil
}

func (s *MemoryWebhookReceiptStore) Confirm(_ context.Context, provider, orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := provider + "\x00" + orderID
	receipt, exists := s.receipts[key]
	if !exists {
		return ErrWebhookReceiptUnavailable
	}
	receipt.accepted = true
	s.receipts[key] = receipt
	return nil
}

func (s *MemoryWebhookReceiptStore) Release(_ context.Context, provider, orderID string) error {
	s.mu.Lock()
	delete(s.receipts, provider+"\x00"+orderID)
	s.mu.Unlock()
	return nil
}

type PostgresWebhookReceiptStore struct{ pool *pgxpool.Pool }

func NewPostgresWebhookReceiptStore(pool *pgxpool.Pool) *PostgresWebhookReceiptStore {
	return &PostgresWebhookReceiptStore{pool: pool}
}

func (s *PostgresWebhookReceiptStore) Claim(ctx context.Context, provider, orderID string, now, expiresAt time.Time) (WebhookReceiptClaim, error) {
	if _, err := s.pool.Exec(ctx, `DELETE FROM webhook_receipts WHERE expires_at <= $1`, now); err != nil {
		return "", ErrWebhookReceiptUnavailable
	}
	var state string
	err := s.pool.QueryRow(ctx, `INSERT INTO webhook_receipts (provider_name,order_id,status,accepted_at,expires_at) VALUES ($1,$2,'pending',$3,$4) ON CONFLICT (provider_name,order_id) DO UPDATE SET status='pending',accepted_at=EXCLUDED.accepted_at,expires_at=EXCLUDED.expires_at WHERE webhook_receipts.expires_at <= $3 RETURNING status`, provider, orderID, now, expiresAt).Scan(&state)
	if err == nil {
		return WebhookReceiptClaimed, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", ErrWebhookReceiptUnavailable
	}
	if err := s.pool.QueryRow(ctx, `SELECT status FROM webhook_receipts WHERE provider_name=$1 AND order_id=$2`, provider, orderID).Scan(&state); err != nil {
		return "", ErrWebhookReceiptUnavailable
	}
	if state == "accepted" {
		return WebhookReceiptDuplicate, nil
	}
	return WebhookReceiptInProgress, nil
}

func (s *PostgresWebhookReceiptStore) Confirm(ctx context.Context, provider, orderID string) error {
	result, err := s.pool.Exec(ctx, `UPDATE webhook_receipts SET status='accepted' WHERE provider_name=$1 AND order_id=$2 AND status='pending'`, provider, orderID)
	if err != nil || result.RowsAffected() != 1 {
		return ErrWebhookReceiptUnavailable
	}
	return nil
}

func (s *PostgresWebhookReceiptStore) Release(ctx context.Context, provider, orderID string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM webhook_receipts WHERE provider_name=$1 AND order_id=$2`, provider, orderID)
	if err != nil {
		return ErrWebhookReceiptUnavailable
	}
	return nil
}

var _ WebhookReceiptStore = (*MemoryWebhookReceiptStore)(nil)
var _ WebhookReceiptStore = (*PostgresWebhookReceiptStore)(nil)
