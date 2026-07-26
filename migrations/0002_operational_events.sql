CREATE TABLE IF NOT EXISTS operational_events (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    order_id TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    attempt INTEGER NOT NULL DEFAULT 0,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS operational_events_occurred_at_idx
    ON operational_events (occurred_at DESC);

CREATE INDEX IF NOT EXISTS operational_events_order_id_idx
    ON operational_events (order_id);
