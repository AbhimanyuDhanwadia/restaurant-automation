# Running Restaurant Automation

## Requirements

| Tool | Version |
|---|---|
| Node.js | 24 (see `.nvmrc`) |
| npm | Ships with Node |
| Go | 1.26+ (for the backend API) |
| Docker & Docker Compose | Latest stable (for the containerised stack) |
| Supabase project | Any plan — needed for authentication |

The repository includes `.nvmrc` for Node version managers. With `nvm`, run:

```bash
nvm install
nvm use
```

## Install

Clone the repository and install the locked dependencies:

```bash
git clone https://github.com/AbhimanyuDhanwadia/restaurant-automation.git
cd restaurant-automation
npm ci
```

## Configure Supabase

Copy the example environment file:

```bash
cp .env.example .env.local
```

Set these values in `.env.local`:

```text
VITE_SUPABASE_URL=https://your-project.supabase.co
VITE_SUPABASE_ANON_KEY=your-publishable-anon-key
```

The URL may also include `/rest/v1/`; the app normalises that suffix for Auth. Use only the publishable anon key in this frontend. Never add a Supabase `service_role` key to `.env.local`, GitHub secrets, or source control.

Create at least one user in Supabase under **Authentication → Users**. The app currently uses email / password sign-in.

## Run the Frontend Locally

Start the Vite development server:

```bash
npm run dev
```

Open the URL printed by Vite, normally `http://localhost:5173`.

### Available Routes

| Path | Page |
|---|---|
| `/` | Operations dashboard |
| `/orders` | Order intake & kitchen queue |
| `/tables` | Table & reservation management |
| `/inventory` | Inventory alerts & purchasing |
| `/staff` | Staff roster & shift handoff |
| `/alerts` | Operational alerts |
| `/operations` | Live operations workspace |
| `/kitchen` | Kitchen ticket queue |
| `/analytics` | Restaurant analytics |
| `/settings` | Restaurant settings |
| `/automation` | Automation health, queues, printers, and events |

All routes are protected by Supabase Auth. You will be redirected to the login screen until valid Supabase credentials are configured.

Kitchen tickets and settings controls remain mock-data interactions. The Analytics page polls the operational reporting endpoint and transparently uses demo values only when the API is unavailable.

## Run the Go API

From the repository root, start the API directly:

```bash
go run ./cmd/server
```

The API listens on `http://localhost:8080`. Check its foundation endpoints:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

Submit a mock order to the Phase 3 automation pipeline and inspect its state:

```bash
curl -X POST http://localhost:8080/api/v1/automation/orders \
  -H 'Content-Type: application/json' \
  -d '{"order_id":"ORD-1842"}'
curl http://localhost:8080/api/v1/automation/queue
curl http://localhost:8080/api/v1/automation/events
curl http://localhost:8080/api/v1/integrations
curl http://localhost:8080/api/v1/printers
curl http://localhost:8080/api/v1/analytics/overview
```

The Phase 3-5 engine is intentionally in-memory. The integrations endpoint reports the local mock provider, and the printers endpoint reports the local mock kitchen printer. To queue a ticket:

```bash
curl -X POST http://localhost:8080/api/v1/printers/tickets \
  -H 'Content-Type: application/json' \
  -d '{"order_id":"ORD-1842","destination":"kitchen","lines":[{"text":"Paneer Tikka","quantity":2}]}'
curl -X POST http://localhost:8080/api/v1/printers/tickets/ORD-1842/reprint
```

The printer manager renders ESC/POS bytes and exercises retry/reconnect behavior through the mock driver; it does not contact physical USB or network printers yet. Restarting the API clears its events, queues, provider state, and printer history.

### Environment Variables (Backend)

These can be set in `.env` or exported in your shell:

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | API listen port |
| `SHUTDOWN_TIMEOUT` | `30s` | Graceful shutdown duration |
| `CORS_ORIGINS` | `http://localhost:5173,http://localhost:5174` | Allowed CORS origins |
| `LOG_LEVEL` | `info` | Zerolog level |
| `LOG_PRETTY` | `true` | Human-readable console logs |
| `DATABASE_URL` | — | PostgreSQL connection string |
| `DB_MAX_CONNS` | `25` | Connection pool max size |
| `DB_MIN_CONNS` | `5` | Connection pool min size |
| `SUPABASE_JWT_SECRET` | — | JWT verification for auth middleware |

## Run With Docker Compose

Docker Compose starts the Go API and PostgreSQL together:

```bash
docker compose -f docker/docker-compose.yml up --build
```

The API is available at `http://localhost:8080` and PostgreSQL at `localhost:5432`. Stop the stack with:

```bash
docker compose -f docker/docker-compose.yml down
```

Add `-v` to the `down` command only when you intentionally want to delete the local PostgreSQL volume.

## Verify Changes

Run the same checks used by GitHub Actions:

```bash
npm run lint          # ESLint (TypeScript + React rules)
npm run build         # tsc type-check + Vite production build
go test ./...         # Go unit tests
go vet ./...          # Go static analysis
```

To preview the production bundle locally:

```bash
npm run preview
```

## Deploy

GitHub Pages deployment is handled by `.github/workflows/deploy.yml`.

1. Add repository secrets named `VITE_SUPABASE_URL` and `VITE_SUPABASE_ANON_KEY`.
2. Set GitHub Pages source to **GitHub Actions**.
3. Add the Pages URL to Supabase Auth redirect URLs:
   `https://abhimanyudhanwadia.github.io/restaurant-automation/`
4. Push to `main`, or manually run **Deploy to GitHub Pages** from the Actions tab.

The deployed site uses the `/restaurant-automation/` base path. Local development uses `/`.

## Automation Dashboard

Open `/automation` after signing in to inspect integrations, worker health, queue metrics, printer status, and the latest event stream. The page polls the Go API every five seconds. Set `VITE_API_URL` when the API is hosted somewhere other than `http://localhost:8080`; otherwise it uses local demo values when the API is unavailable.

## Troubleshooting

| Symptom | Fix |
|---|---|
| **Auth setup screen appears** | Confirm both `VITE_SUPABASE_URL` and `VITE_SUPABASE_ANON_KEY` exist in `.env.local`, then restart Vite. |
| **Invalid login** | Confirm the user exists in Supabase **Authentication → Users** and use the correct email / password. |
| **Pages shows setup screen** | Confirm both repository secrets are set and re-run the deployment workflow. |
| **Stale local data** | Operational state is stored in browser `localStorage`. Clear site storage in your browser to reset to demo data. |
| **Port 8080 in use** | Set `PORT=9090` (or another free port) in `.env` before starting the Go API. |
| **Docker Compose fails to start** | Ensure Docker is running and port 5432 is free. Check `docker compose logs` for details. |
| **Go module errors** | Run `go mod tidy` then retry. |
