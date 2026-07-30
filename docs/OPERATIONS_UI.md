# Operations Dashboard API Connection

The Active Orders panel on Live Operations uses the durable Orders API.

## Live data

- Summary cards derive counts from all durable orders: active, queued, preparing, ready, and delivered.
- The active-order list includes only `received`, `preparing`, and `ready` orders.
- Topbar search filters the API-backed active-order list by order ID, channel, item name, notes, or status.
- The global **New order** action opens the durable Orders workspace, where staff create the order through the API.
- The Next Tasks panel uses the durable Staff task feed and completes tasks through the API.

Alerts remain a separate local module. Tables, inventory, staff roster, and staff tasks have their own durable workspaces.
