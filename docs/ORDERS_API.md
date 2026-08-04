# Orders API

`POST /api/v1/orders` creates an order with `channel`, optional `notes`, at least one item, and optional `total_minor` plus a three-letter `currency`. Totals are stored in minor units and must be non-negative. `GET /api/v1/orders` lists orders. `PATCH /api/v1/orders/{orderID}/status` accepts received, preparing, ready, delivered, or cancelled.

When PostgreSQL is configured, orders and items are stored durably; otherwise the API uses an in-memory repository for local development.
