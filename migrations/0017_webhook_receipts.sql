CREATE TABLE IF NOT EXISTS webhook_receipts (
  provider_name TEXT NOT NULL,
  order_id TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pending', 'accepted')),
  accepted_at TIMESTAMPTZ NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (provider_name, order_id)
);

CREATE INDEX IF NOT EXISTS webhook_receipts_expires_at_index ON webhook_receipts (expires_at);
