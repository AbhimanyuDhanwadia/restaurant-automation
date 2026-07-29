CREATE TABLE IF NOT EXISTS inventory_items (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  category TEXT NOT NULL,
  on_hand DOUBLE PRECISION NOT NULL CHECK (on_hand >= 0),
  par_level DOUBLE PRECISION NOT NULL CHECK (par_level >= 0),
  unit TEXT NOT NULL,
  supplier TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL CHECK (status IN ('in_stock', 'low_stock', 'on_order')),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
