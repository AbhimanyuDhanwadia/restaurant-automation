# Database API

`GET /api/v1/admin/database` is a read-only Administration endpoint requiring `database.view`.

When PostgreSQL is configured, it returns `status: "available"` and the applied versions from `schema_migrations`, newest first. When PostgreSQL is intentionally absent for local development, it returns `status: "not_configured"` with an empty migration list.

It never returns connection strings, credentials, tables, or arbitrary query access. No migration is required because the `schema_migrations` ledger already exists.
