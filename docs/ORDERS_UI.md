# Orders UI API Connection

The Orders workspace reads and writes through the Go API rather than the local mock store.

## Configuration

Set the frontend API address before starting Vite:

```bash
VITE_API_URL=http://localhost:8080
```

The default is `http://localhost:8080`. When Supabase authentication is configured and a user has signed in, the frontend includes that user's access token in the API request. This supports deployments with `AUTH_REQUIRED=true`.

## Supported operations

- Load all durable orders from `GET /api/v1/orders`.
- Create an order with channel, item, quantity, and optional notes through `POST /api/v1/orders`.
- Update an order to `received`, `preparing`, `ready`, `delivered`, or `cancelled` through `PATCH /api/v1/orders/{orderID}/status`.

The Orders screen shows a visible error and retry control when the API cannot be reached. The full REST contract is in [ORDERS_API.md](ORDERS_API.md).
