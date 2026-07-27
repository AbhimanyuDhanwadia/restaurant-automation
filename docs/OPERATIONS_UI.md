# Operations Dashboard API Connection

The Active Orders panel on Live Operations uses the durable Orders API.

## Live data

- Summary cards derive counts from all durable orders: active, queued, preparing, ready, and delivered.
- The active-order list includes only `received`, `preparing`, and `ready` orders.
- Topbar search filters the API-backed active-order list by order ID, channel, item name, notes, or status.
- The global **New order** action opens the durable Orders workspace, where staff create the order through the API.

Alerts, tables, inventory, and staff tasks remain separate modules and are not represented as API-backed data in this dashboard milestone.
