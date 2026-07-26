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
