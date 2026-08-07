# Orders API

`POST /api/v1/orders` creates an order with `channel`, optional `notes`, at least one item, optional `total_minor` plus a three-letter `currency`, and optional `delivery_partner`. Totals are stored in minor units and must be non-negative. A successful creation also creates one kitchen-destination print job using the persisted item list, then submits the order to the automation engine. When an order with a delivery partner first reaches `delivered`, the API records `delivered_at`. `GET /api/v1/orders` lists orders. `PATCH /api/v1/orders/{orderID}/status` accepts received, preparing, ready, delivered, or cancelled.

When PostgreSQL is configured, orders and items are stored durably; otherwise the API uses an in-memory repository for local development.

The order is still accepted when dispatching an already-created kitchen job fails, such as when no kitchen printer route is configured. The job is stored with `failed` status and can be recovered through the print-queue reprint endpoint. If the durable print job itself cannot be created, the request returns `503 Service Unavailable` after the order record has been created; clients should check the order list before retrying.

Changing a status emits one lifecycle event after the durable order update: `preparing` emits `kitchen.accepted`, `ready` emits `order.ready`, `delivered` emits `order.delivered`, and `cancelled` emits `order.cancelled`. Repeating the same status does not create another event. The order record remains authoritative if the in-memory event engine is unavailable, so clients must not retry a successful status response solely to recover event processing.
