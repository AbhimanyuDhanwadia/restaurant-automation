# Sales UI

The Sales workspace refreshes authenticated `GET /api/v1/analytics/overview` data every 15 seconds. It displays recorded sales, average ticket, total coverage, excluded eligible orders, and the reporting currency.

Sales include only non-cancelled durable orders with `total_minor` and `currency`. Orders without totals are excluded and remain visible in the coverage count. When recorded orders use multiple currencies, revenue and average ticket are deliberately unavailable rather than combined incorrectly.

Order totals are entered as decimal currency values in the Orders workspace and are stored as integer minor units. This avoids floating-point storage errors.
