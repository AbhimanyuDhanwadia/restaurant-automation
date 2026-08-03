CREATE TABLE IF NOT EXISTS print_jobs (
  id UUID PRIMARY KEY,
  order_id TEXT NOT NULL,
  destination TEXT NOT NULL,
  lines JSONB NOT NULL,
  reprint BOOLEAN NOT NULL DEFAULT FALSE,
  status TEXT NOT NULL CHECK (status IN ('queued', 'printing', 'printed', 'failed')),
  attempts INTEGER NOT NULL DEFAULT 0,
  last_error TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS print_jobs_created_at_idx ON print_jobs (created_at DESC);
CREATE INDEX IF NOT EXISTS print_jobs_order_id_idx ON print_jobs (order_id, created_at DESC);
