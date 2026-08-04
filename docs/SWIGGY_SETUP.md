# Swiggy Setup

Restaurant Automation is ready to receive signed Swiggy webhooks once Swiggy
provides the restaurant/POS integration credentials and contract.

## Configure

Set the signing secret supplied during Swiggy onboarding:

```text
SWIGGY_WEBHOOK_SECRET=replace-with-the-swiggy-signing-secret
```

Restart the API. The provider registry will expose `swiggy` as a connected
integration and accept deliveries at:

```text
POST /api/v1/webhooks/swiggy
```

The endpoint is public only for delivery from Swiggy. It verifies the raw
payload with `X-Webhook-Signature` using HMAC-SHA256 before placing the order
on the provider collector queue.

For the current neutral order envelope, an accepted `order_id` is remembered
for 15 minutes. A duplicate signed delivery returns HTTP `202` with
`{"status":"duplicate"}` and is not queued again. This prevents duplicate
automation and print work during normal webhook retries. With PostgreSQL
configured, receipts survive API restarts and are shared across instances;
without it, the local in-memory fallback resets on restart. This is not a
substitute for a provider-specified timestamp or nonce scheme.

If another API instance is currently processing the same order, the endpoint
returns a retryable HTTP `503` rather than acknowledging it as a duplicate.

## Current Payload Boundary

Until Swiggy supplies the POS webhook schema, the receiver expects this
integration-neutral envelope:

```json
{
  "order_id": "SWIGGY-ORDER-ID",
  "payload": {}
}
```

The `payload` object is retained for the later Swiggy-specific normalization
mapping. Do not send production traffic until the following are confirmed with
the Swiggy partner team:

1. Exact webhook URL registration and retry behavior.
2. Signature header name, signing algorithm, and timestamp/nonce replay protection.
3. Order payload schema, including line items and order status events.
4. Required acknowledgement and status-update APIs.
5. Sandbox credentials and production rollout process.

The public [Swiggy Developer Portal](https://developers.swiggy.com/) and
[Swiggy Builders Club developer quickstart](https://mcp.swiggy.com/builders/docs/start/developer/)
describe the available developer platform, but they do not publish the
restaurant/POS webhook payload contract required for a complete adapter.
