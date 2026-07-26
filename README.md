# Restaurant Automation

A web operations console for restaurants to coordinate live orders, tables, kitchen load, inventory alerts, and staff handoffs from one screen.

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React 19, TypeScript, Tailwind CSS 4, Vite |
| Routing | TanStack Router (code-based, lazy-loaded) |
| Data / State | TanStack React Query, Zustand (localStorage-persisted) |
| Auth | Supabase Auth (email / password) |
| Backend API | Go 1.26, Chi v5, Zerolog, Viper |
| Database | PostgreSQL 16 (Docker Compose) |
| CI / CD | GitHub Actions → GitHub Pages |
| Icons | Lucide React |

## Features

- **Operations Dashboard** — summary stat cards for open orders, kitchen load, table occupancy, and online staff.
- **Orders** — live order list with status cycling (Preparing → Ready), channel tags (Dine-in, Delivery, Takeaway), and a new-order form.
- **Tables** — floor-plan directory with seat counts, guest names, reservation times, and a seating lifecycle (Available → Reserved → Seated → Needs Check).
- **Inventory** — stock tracker with par levels, supplier info, category filtering, and one-click receive action.
- **Staff** — shift roster, shift-task queue with completion, and a manager handoff note that persists across reloads.
- **Alerts** — severity-tagged operational alerts (High / Medium / Low) with acknowledge-to-dismiss.
- **Auth** — Supabase-backed login gate; every route is protected by `AuthGuard`.
- **Local-first persistence** — all domain stores use Zustand `persist` middleware with automatic legacy-format migration.
- **Global search** — top-bar search filters the active page's data in real time.

## v2 Phase Status

- Phase 1: Foundation complete. Go API, React frontend, PostgreSQL Compose setup, Supabase Auth, logging, CI, and environment configuration are in place.
- Phase 2: Restaurant UI complete on `v2`. Dashboard, Orders, Kitchen, Tables, Inventory, Staff, Alerts, Analytics, and Settings use mock data and local UI state.
- Phase 3: Automation engine complete on `v2`. The Go API now provides an in-memory event bus, bounded worker queue, retry policy, scheduler, order pipeline, and queue/event inspection endpoints.
- Phase 4+: Not started. Provider integrations, persistent storage, printer infrastructure, automation dashboard, analytics persistence, and AI remain separate future phases.

## Project Structure

```
.
├── cmd/server/           # Go API entrypoint
├── docker/               # Dockerfile & docker-compose.yml
├── internal/
│   ├── api/              # Chi router, handlers, middleware
│   ├── automation/       # Phase 3 event bus, queues, retries, workers
│   ├── config/           # Viper-based configuration
│   └── logger/           # Zerolog setup
├── migrations/           # SQL migration files
├── src/
│   ├── components/
│   │   ├── layout/       # AppShell, Sidebar, Topbar
│   │   └── ui/           # Panel, StatusPill, DirectoryRow, EmptyState
│   ├── data/             # Seed / demo data
│   ├── features/auth/    # AuthGuard, LoginPage
│   ├── lib/              # Supabase client, query client, utils, storage migration
│   ├── pages/            # Operations, Orders, Kitchen, Tables, Inventory, Staff, Alerts, Analytics, Settings
│   ├── stores/           # Zustand stores (orders, tables, inventory, staff, alerts, ui)
│   ├── types/            # Domain type definitions
│   ├── router.tsx        # TanStack Router route tree
│   ├── main.tsx          # React entry
│   └── styles.css        # Global styles / Tailwind entry
├── .github/workflows/    # CI (ci.yml) and deploy (deploy.yml)
├── .env.example          # Environment variable template
└── package.json
```

## Milestones

1. ~~App foundation and operations dashboard.~~
2. ~~Order intake and kitchen queue workflow.~~
3. ~~Table, reservation, and guest status management.~~
4. ~~Inventory alerts and purchasing workflow.~~
5. ~~Staff task automation and shift handoff views.~~
6. Persistence, authentication, deployment, and production hardening.

## Local Development

```bash
npm ci
npm run dev
```

See [RUNNING.md](RUNNING.md) for Node.js requirements, Supabase setup, Docker Compose, verification commands, deployment, and troubleshooting.

## Verification

Every push to `main` or `v2`, and every pull request targeting either branch, runs the CI workflow. It installs the locked dependencies, runs frontend lint / build checks, and runs `go test ./...` plus `go vet ./...`.

## Backend API

The Go API starts on port `8080` and exposes:

| Endpoint | Purpose |
|---|---|
| `GET /health` | Liveness probe |
| `GET /ready` | Readiness probe |
| `POST /api/v1/automation/orders` | Submit an order to the Phase 3 pipeline |
| `GET /api/v1/automation/events` | Inspect the in-memory event stream |
| `GET /api/v1/automation/queue` | Inspect queue depth and worker metrics |
| `/api/v1/*` | Versioned API namespace |

PostgreSQL and the API can be started together with Docker Compose. Configuration is loaded from environment variables or a `.env` file; secrets are never committed.

Phase 3 intentionally keeps order processing in memory. It demonstrates collection, normalization, validation, persistence boundary, event publication, and queueing without claiming database or printer delivery. State resets when the Go API restarts; PostgreSQL persistence, provider adapters, and ESC/POS infrastructure are later phases.

## Deployment

The **Deploy to GitHub Pages** workflow builds and publishes `main` to GitHub Pages with Node.js 24. Enable GitHub Pages for the repository with the **GitHub Actions** source, add the Supabase repository secrets, then trigger the workflow from the Actions tab or push to `main`.

## License

[MIT](LICENSE) © Abhimanyu Dhanwadia
