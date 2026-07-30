CREATE TABLE IF NOT EXISTS staff_members (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  role TEXT NOT NULL,
  station TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('on_shift', 'on_break', 'off_shift')),
  handoff TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS shift_tasks (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  owner TEXT NOT NULL,
  due_label TEXT NOT NULL DEFAULT '',
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS shift_handoffs (
  id TEXT PRIMARY KEY,
  note TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL
);
