CREATE TABLE IF NOT EXISTS backup_records (
  id UUID PRIMARY KEY,
  target TEXT NOT NULL,
  size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
  status TEXT NOT NULL CHECK (status IN ('recorded', 'verified', 'failed')),
  completed_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS backup_records_completed_at_index ON backup_records (completed_at DESC);
