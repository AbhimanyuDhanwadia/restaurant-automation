# Webhook Receipts

Migration `0017_webhook_receipts.sql` adds a durable idempotency store keyed by provider and order ID. Before a signed webhook order is queued, the provider atomically claims a pending receipt for 15 minutes and marks it accepted only after enqueueing. A second API instance or a restarted process sees an accepted receipt and returns an acknowledged duplicate response instead of producing another order. A concurrent pending receipt returns a retryable service response rather than a false duplicate acknowledgement.

If the provider queue is full, the receipt is released so the delivery provider can retry later. Expired receipts are deleted during subsequent claim attempts. Local development without PostgreSQL uses the same interface with an in-memory store, which resets at API restart.

This protects the current generic order envelope. Timestamp, nonce, and event-type semantics remain dependent on the provider's confirmed production webhook contract.
