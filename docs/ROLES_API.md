# Roles API

Migration `0012_roles.sql` creates durable `app_roles` records and seeds Operator, Manager, and Administrator system templates. A role has a name, description, normalized allow-listed permissions, and a `system` marker.

`GET /api/v1/admin/roles` lists roles. `POST /api/v1/admin/roles` creates a custom role. Role names must be unique and at least one supported permission is required. Observed identities can be assigned a role through the Users administration workspace; see `USERS_API.md`. Role permissions are enforced for protected API capabilities; see `AUTHORIZATION.md`.
