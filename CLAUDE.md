# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

Full-stack 3D printing service. Backend is a Go REST API (`server/`), frontend is a TypeScript/React app using TanStack Start (`client/`), database is MySQL. PrusaSlicer is invoked server-side to slice uploaded STL/3MF files and calculate print costs.

## Development commands

### Full stack (Docker)

```bash
docker compose up          # MySQL + API (air hot reload) + Client on ports 3306/8080/3000
```

### Server (Go)

```bash
cd server
make build                 # compile to bin/server
make run                   # build + run
make test                  # go test -v ./...
go test -v ./services/filaments  # run a single package's tests

make migrate-up            # apply pending migrations
make migrate-down          # rollback latest migration
make migration NAME=<name> # create new up/down migration pair
make seed                  # load seeds from cmd/migrate/seeds/
```

### Client (TypeScript/React)

```bash
cd client
npm run dev                # Vite dev server on :3000
npm run build              # production build → .output/
npm test                   # vitest run
```

## Architecture

### Server

Entry point: `cmd/main.go` → initializes DB + starts HTTP server on port 8080.

**Service pattern** — each service under `services/<name>/` has:
- `routes.go` — HTTP handlers registered on a Gorilla Mux router
- `repository.go` — direct SQL queries (no ORM)

Services are wired in `cmd/api/api.go` which calls `NewRouter()` on each and passes the shared `*mux.Router`.

**Current services:** `health`, `filaments`, `orders`, `prints`, `quotes` (file upload + PrusaSlicer slicing), `jobs` (WIP — routes defined, handlers empty).

**Shared types** — repository interfaces and domain structs live in `types/types.go`.

**Config** — loaded from `.env` via `config/env.go` (godotenv). See `.env.example` for required vars: `DB_*`, `PORT`, `UPLOAD_PATH`, `PRUSA_SLICER_PATH`, `PRUSA_SLICER_CONFIG`.

**Migrations** — `cmd/migrate/migrations/` numbered pairs (`N_name.up.sql` / `N_name.down.sql`), run via golang-migrate through the Makefile.

**UUIDs** — stored as binary in MySQL; queries use `UUID_TO_BIN()` / `BIN_TO_UUID()`.

**Error responses** — use `utils.WriteError(w, status, err)` for consistent JSON errors.

### Client

File-based routing via TanStack Router under `src/routes/`.

- `src/lib/api.ts` — Axios instance configured with `VITE_API_URL`
- `src/lib/queries.ts` — TanStack Query hooks used across routes
- `src/components/` — shared UI components
- Styling via Tailwind CSS 4

### Docker

- **Dev:** `docker-compose.yml` — API uses `Dockerfile.dev` with Air hot reload; client uses `Dockerfile.dev` with Vite HMR; source directories are bind-mounted.
- **Prod:** `docker-compose.prod.yml` — multi-stage build; Go binary compiled statically; tests run during build.
