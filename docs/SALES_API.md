# Sales Data API

Migration `0010_order_totals.sql` adds nullable `total_minor` and `currency` columns to `orders`. `total_minor` is a non-negative integer in the currency's minor unit, and `currency` is a three-letter uppercase code when a total is present.

`POST /api/v1/orders` accepts optional `total_minor` and `currency`; both must be supplied together when recording a total. `GET /api/v1/orders` returns these fields. `GET /api/v1/analytics/overview` returns `sales` and `average_ticket` with `value`, `available`, `currency`, `included_orders`, and `excluded_orders`.
