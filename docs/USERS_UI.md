# Users UI

The Administration > Users workspace lists identities that have already made an authenticated API request. Administrators can assign any role from the durable role catalog and can refresh the directory without leaving the workspace.

The page intentionally has no create, delete, password-reset, or Supabase-user listing control. Account lifecycle remains in Supabase Auth. An API response of `administrator access required` means the signed-in email is absent from the API's `ADMIN_EMAILS` environment variable.
