CREATE TABLE IF NOT EXISTS restaurant_settings (
  id TEXT PRIMARY KEY,
  restaurant_name TEXT NOT NULL,
  timezone TEXT NOT NULL,
  operational_alerts BOOLEAN NOT NULL DEFAULT TRUE,
  auto_advance_tickets BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at TIMESTAMPTZ NOT NULL
);
