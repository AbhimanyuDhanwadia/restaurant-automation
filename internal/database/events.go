package database

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/rs/zerolog"
)

type EventStore struct{ pool *pgxpool.Pool }

func NewEventStore(pool *pgxpool.Pool) *EventStore { return &EventStore{pool: pool} }

func (s *EventStore) Save(ctx context.Context, event automation.Event) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("marshal event payload: %w", err)
	}
	_, err = s.pool.Exec(ctx, `INSERT INTO operational_events (id, event_type, order_id, occurred_at, attempt, payload) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO NOTHING`, event.ID, event.Type, event.OrderID, event.CreatedAt, event.Attempt, payload)
	if err != nil {
		return fmt.Errorf("insert operational event: %w", err)
	}
	return nil
}

func (s *EventStore) List(ctx context.Context, limit int) ([]automation.Event, error) {
	if limit < 1 {
		limit = 100
	}
	if limit > 1_000 {
		limit = 1_000
	}
	rows, err := s.pool.Query(ctx, `SELECT id, event_type, order_id, occurred_at, attempt, payload FROM operational_events ORDER BY occurred_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("list operational events: %w", err)
	}
	defer rows.Close()

	result := make([]automation.Event, 0)
	for rows.Next() {
		var event automation.Event
		var eventType string
		var payload []byte
		if err := rows.Scan(&event.ID, &eventType, &event.OrderID, &event.CreatedAt, &event.Attempt, &payload); err != nil {
			return nil, fmt.Errorf("scan operational event: %w", err)
		}
		event.Type = automation.EventType(eventType)
		if len(payload) > 0 {
			if err := json.Unmarshal(payload, &event.Payload); err != nil {
				return nil, fmt.Errorf("unmarshal event payload: %w", err)
			}
		}
		result = append(result, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate operational events: %w", err)
	}
	return result, nil
}

// EventPersister observes the event bus without becoming a dependency of the
// automation engine. A failed write is logged and does not block order flow.
type EventPersister struct {
	store *EventStore
	log   zerolog.Logger
	stop  func()
	wg    sync.WaitGroup
}

func NewEventPersister(store *EventStore, log zerolog.Logger) *EventPersister {
	return &EventPersister{store: store, log: log}
}

func (p *EventPersister) Start(ctx context.Context, engine *automation.Engine) {
	events, stop := engine.Subscribe(256)
	p.stop = stop
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-events:
				if !ok {
					return
				}
				if err := p.store.Save(ctx, event); err != nil {
					p.log.Error().Err(err).Str("event_id", event.ID).Msg("persist operational event")
				}
			}
		}
	}()
}
func (p *EventPersister) Close() {
	if p.stop != nil {
		p.stop()
	}
	p.wg.Wait()
}
