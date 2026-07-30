CREATE TABLE IF NOT EXISTS operational_alerts (
  id UUID PRIMARY KEY,
  label TEXT NOT NULL,
  detail TEXT NOT NULL,
  severity TEXT NOT NULL CHECK (severity IN ('high', 'medium', 'low')),
  acknowledged_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL
);
