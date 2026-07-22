# Running Restaurant Automation

## Requirements

- Node.js 24
- npm
- A Supabase project for authentication

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

The URL may also include `/rest/v1/`; the app normalizes that suffix for Auth. Use only the publishable anon key in this frontend. Never add a Supabase `service_role` key to `.env.local`, GitHub secrets, or source control.

Create at least one user in Supabase under **Authentication → Users**. The app currently uses email/password sign-in.

## Run Locally

Start the Vite development server:

```bash
npm run dev
```

Open the URL printed by Vite, normally `http://localhost:5173`.

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

## Run With Docker Compose

Docker Compose starts the Go API and PostgreSQL together:

```bash
docker compose -f docker/docker-compose.yml up --build
```

The API is available at `http://localhost:8080` and PostgreSQL is available at `localhost:5432`. Stop the stack with:

```bash
docker compose -f docker/docker-compose.yml down
```

Add `-v` to the `down` command only when you intentionally want to delete the local PostgreSQL volume.

## Verify Changes

Run the same checks used by GitHub Actions:

```bash
npm run lint
npm run build
go test ./...
go vet ./...
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

## Troubleshooting

- **Authentication setup screen:** confirm both `VITE_SUPABASE_URL` and `VITE_SUPABASE_ANON_KEY` exist in `.env.local`, then restart Vite.
- **Invalid login:** confirm the user exists in Supabase Authentication and use the correct email/password.
- **Pages setup screen:** confirm both repository secrets are set and rerun the deployment workflow.
- **Stale local data:** operational changes are stored in browser local storage. Clear the site storage in browser settings to return to the initial demo data.
