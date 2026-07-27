# Tables API

All endpoints are under `/api/v1` and use the standard API authentication middleware.

- `GET /tables` lists tables by ID.
- `POST /tables` creates a table with `id`, positive `seats`, optional `guest`, optional `reservation`, and optional status.
- `PATCH /tables/{tableID}/status` updates a table status.

Valid statuses are `available`, `seated`, `reserved`, and `needs_check`.

The `restaurant_tables` data is persisted when PostgreSQL is configured; local development without a database uses the in-memory repository.
