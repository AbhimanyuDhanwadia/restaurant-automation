# Delivery Data API

Migration `0011_delivery_tracking.sql` adds nullable `delivery_partner` and `delivered_at` columns to `orders`. `POST /api/v1/orders` accepts an optional `delivery_partner`. When an order with this value is first updated to `delivered`, the API writes `delivered_at`; later updates retain the original completion timestamp.

`GET /api/v1/analytics/overview` returns `delivery_time` with `value` in minutes, `available`, `delivered_orders`, and `active_orders`. It includes only non-cancelled orders with a delivery partner.
