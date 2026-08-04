UPDATE app_roles
SET permissions = (
  SELECT jsonb_agg(DISTINCT permission ORDER BY permission)
  FROM jsonb_array_elements_text(app_roles.permissions || '["backups.manage","users.manage"]'::jsonb) AS expanded(permission)
), updated_at = NOW()
WHERE id = '00000000-0000-0000-0000-000000000003';
