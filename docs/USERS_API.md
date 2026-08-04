# Users API

Migration `0014_app_users.sql` creates the local `app_users` directory. It stores a verified authentication subject, observed email address, role assignment, and timestamps. It does not copy or manage Supabase Auth accounts.

Every authenticated `/api/v1` request records the JWT subject and email. The first observation receives the `Operator` role, except emails listed in `ADMIN_EMAILS`, which receive the `Administrator` role. Later observations never overwrite a role chosen by an administrator.

`GET /api/v1/admin/users` lists observed identities. `PATCH /api/v1/admin/users/{subject}/role` accepts `{ "role_id": "<role UUID>" }` and changes a user's local application role. Both routes require a valid Supabase JWT and an email in `ADMIN_EMAILS`.

Set `AUTH_REQUIRED=true`, `SUPABASE_JWT_SECRET`, and `ADMIN_EMAILS=manager@example.com` on the API to enable the directory in a deployed environment. `ADMIN_EMAILS` is only a bootstrap allow-list for this administration endpoint; broader role-based route authorization is a later milestone. No Supabase service-role key is used or needed.
