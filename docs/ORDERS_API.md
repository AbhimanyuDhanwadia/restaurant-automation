# Orders API

`POST /api/v1/orders` creates an order with `channel`, optional `notes`, at least one item, optional `total_minor` plus a three-letter `currency`, and optional `delivery_partner`. Totals are stored in minor units and must be non-negative. When an order with a delivery partner first reaches `delivered`, the API records `delivered_at`. `GET /api/v1/orders` lists orders. `PATCH /api/v1/orders/{orderID}/status` accepts received, preparing, ready, delivered, or cancelled.

When PostgreSQL is configured, orders and items are stored durably; otherwise the API uses an in-memory repository for local development.
