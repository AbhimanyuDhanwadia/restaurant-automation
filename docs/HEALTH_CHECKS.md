# Health Checks

`GET /health` reports that the HTTP server is accepting requests.

`GET /ready` additionally pings PostgreSQL when `DATABASE_URL` is configured.
It returns `503` when that configured database is unavailable, allowing Docker,
load balancers, and deployment platforms to avoid routing traffic to an API
that cannot persist operational events.
