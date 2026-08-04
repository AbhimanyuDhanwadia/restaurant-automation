CREATE TABLE IF NOT EXISTS app_roles (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
  system BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS app_roles_name_lower_unique ON app_roles (LOWER(name));

INSERT INTO app_roles (id, name, description, permissions, system, created_at, updated_at) VALUES
  ('00000000-0000-0000-0000-000000000001', 'Operator', 'Runs day-to-day restaurant operations.', '["operations.view","orders.manage","kitchen.manage","inventory.manage"]'::jsonb, TRUE, NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000002', 'Manager', 'Oversees operations, staffing, and analytics.', '["operations.view","orders.manage","kitchen.manage","inventory.manage","staff.manage","analytics.view","settings.manage"]'::jsonb, TRUE, NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000003', 'Administrator', 'Configures automation and administration services.', '["operations.view","orders.manage","kitchen.manage","inventory.manage","staff.manage","analytics.view","automation.view","integrations.manage","printers.manage","settings.manage","audit.view","database.view","roles.manage"]'::jsonb, TRUE, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;
