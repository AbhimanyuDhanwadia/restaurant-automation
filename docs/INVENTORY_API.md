# Inventory API

All endpoints are under `/api/v1` and use the standard API authentication middleware.

- `GET /inventory` lists inventory items by name.
- `POST /inventory` creates an item with name, category, non-negative `on_hand`, non-negative `par_level`, unit, and optional supplier or status.
- `PATCH /inventory/{itemID}/stock` sets `on_hand` and automatically derives `in_stock` or `low_stock` from the par level.
- `PATCH /inventory/{itemID}/status` explicitly sets `in_stock`, `low_stock`, or `on_order`.

Inventory data is persisted with PostgreSQL. The in-memory repository supports database-free local development.
