CREATE TABLE IF NOT EXISTS app_users (
  auth_subject TEXT PRIMARY KEY,
  email TEXT NOT NULL,
  role_id UUID NOT NULL REFERENCES app_roles(id),
  created_at TIMESTAMPTZ NOT NULL,
  last_seen_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS app_users_email_index ON app_users (LOWER(email));
CREATE INDEX IF NOT EXISTS app_users_role_id_index ON app_users (role_id);
