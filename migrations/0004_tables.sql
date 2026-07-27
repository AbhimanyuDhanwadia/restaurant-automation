CREATE TABLE IF NOT EXISTS restaurant_tables (
  id TEXT PRIMARY KEY,
  seats INTEGER NOT NULL CHECK (seats > 0),
  guest TEXT NOT NULL DEFAULT '',
  reservation TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL CHECK (status IN ('available', 'seated', 'reserved', 'needs_check')),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
