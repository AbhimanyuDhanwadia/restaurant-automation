# Restaurant Automation

A web operations console for restaurants to coordinate live orders, tables, kitchen load, inventory alerts, and staff handoffs from one screen. It uses React, Vite, Supabase Auth, Go, PostgreSQL, Docker Compose, and GitHub Actions.

The `v2` branch is the active development line for the local-first platform blueprint. Work is delivered phase by phase and pushed after each completed phase.

## Milestones

1. App foundation and operations dashboard.
2. Order intake and kitchen queue workflow.
3. Table, reservation, and guest status management.
4. Inventory alerts and purchasing workflow.
5. Staff task automation and shift handoff views.
6. Persistence, authentication, deployment, and production hardening.

This project is being built in reviewable milestones. Each milestone should leave the repository in a runnable state before moving to the next one.

## Local Development

```bash
npm ci
npm run dev
```

See [RUNNING.md](RUNNING.md) for Node.js requirements, Supabase setup, verification commands, deployment configuration, and troubleshooting.

## Verification

Every push to `main` or `v2`, and every pull request targeting either branch, runs the CI workflow. It installs the locked dependencies, runs frontend lint/build checks, and runs `go test ./...` plus `go vet ./...`.

## Backend Foundation

The Go API starts on port `8080` and currently exposes `/health`, `/ready`, and the versioned `/api/v1` route namespace. PostgreSQL and the API can be started together with Docker Compose. Configuration is loaded from environment variables or a local `.env` file; secrets are never committed.

## Deployment

The `Deploy to GitHub Pages` workflow builds and publishes `main` to GitHub Pages with Node.js 24. Enable GitHub Pages for the repository with the `GitHub Actions` source, add the Supabase repository secrets, then trigger the workflow from the Actions tab or push to `main`.
