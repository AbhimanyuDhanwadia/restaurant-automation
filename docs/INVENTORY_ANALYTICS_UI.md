# Inventory Analytics UI

The Inventory Analytics workspace refreshes authenticated `GET /api/v1/inventory` data every 15 seconds. It reports item-count coverage, low-stock and on-order counts, category-level risk, and individual replenishment shortfalls.

Coverage is a count of items marked in stock, not a summed quantity. Each item is compared only with its own par level, which avoids combining incompatible units such as kilograms, litres, and pieces.

Cost, consumption velocity, purchase orders, and supplier lead-time analytics require future inventory transaction and vendor data. No database migration or API contract was required for this read-only milestone.
