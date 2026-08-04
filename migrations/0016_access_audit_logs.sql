CREATE TABLE IF NOT EXISTS access_audit_logs (
  id UUID PRIMARY KEY,
  actor_subject TEXT NOT NULL,
  actor_email TEXT NOT NULL,
  target_subject TEXT NOT NULL,
  previous_role_id UUID NOT NULL REFERENCES app_roles(id),
  new_role_id UUID NOT NULL REFERENCES app_roles(id),
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS access_audit_logs_created_at_index ON access_audit_logs (created_at DESC);
CREATE INDEX IF NOT EXISTS access_audit_logs_target_subject_index ON access_audit_logs (target_subject);
