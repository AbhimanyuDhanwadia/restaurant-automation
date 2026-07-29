# Inventory UI API Connection

The Inventory workspace uses the shared authenticated frontend API client.

Staff can create inventory items, filter durable stock records, update the on-hand quantity, and set an item to on order. Updating quantity recalculates its in-stock or low-stock state against its stored par level.

The global **New inventory item** action opens this workspace so creation always uses the API.
