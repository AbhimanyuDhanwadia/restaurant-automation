# Authorization

With `AUTH_REQUIRED=true`, every `/api/v1` request is verified against the Supabase JWT secret and recorded in the local application user directory. Authorization then resolves the user's assigned `app_roles` permission list. When `AUTH_REQUIRED=false`, permission checks intentionally allow local development traffic.

The protected API capabilities are:

- `analytics.view`: analytics overview
- `automation.view`: system health, AI insights, and automation engine routes
- `integrations.manage`: integrations
- `printers.manage`: printers, print queue, and ticket actions
- `settings.manage`: settings
- `audit.view`: audit logs
- `database.view`: database status
- `roles.manage`: role catalog
- `users.manage`: user directory and role assignment
- `backups.manage`: backup artifact catalog

Restaurant operations endpoints remain authenticated but are not permission-gated in this rollout. The `ADMIN_EMAILS` setting only assigns the Administrator role on a user's first observed request; it does not bypass subsequent permission checks.

`GET /api/v1/me` returns the verified subject and email together with the local role and current permissions. The authenticated frontend uses this profile to hide protected navigation destinations after it has loaded. The API remains authoritative and returns `403` for a manually entered protected route without the required permission.
