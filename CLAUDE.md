# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

Full-stack 3D printing service. Backend is a Go REST API (`server/`), customer-facing frontend is a TypeScript/React app (`client/`), admin panel is a separate React app (`admin/`), database is MySQL. PrusaSlicer is invoked server-side to slice uploaded STL/3MF files and calculate print costs.

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
make seed                  # load seeds (requires docker compose MySQL running)
```

### Client (TypeScript/React) — customer frontend

```bash
cd client
npm run dev                # Vite dev server on :3000
npm run build              # production build → .output/
npm test                   # vitest run
npm run format             # prettier + eslint fix
```

### Admin panel

```bash
cd admin
npm dev                   # Vite dev server on :3000
npm build
npm test
```

## Architecture

### Server

Entry point: `cmd/main.go` → initializes DB + starts HTTP server on port 8080.

**Service pattern** — each service under `services/<name>/` contains:
- `routes.go` — HTTP handlers registered on a Gorilla Mux router
- `repository.go` — direct SQL queries (no ORM)
- `service.go` — domain logic (present in some services, e.g. `quotes`)

Services are wired in `cmd/api/api.go` which calls `RegisterRoutes()` on each, passing the shared `*mux.Router` prefixed at `/api/v1`.

**Current services:** `health`, `filaments`, `orders`, `prints`, `quotes`, `jobs` (WIP — routes defined, handlers empty).

**Quotes + pricing** (`services/quotes/service.go`):
- `SliceFilament()` invokes PrusaSlicer CLI, writes gcode to a temp file, then parses `; filament used [cm3] =` from the gcode
- Converts cm³ → grams using a hardcoded density of 1.25 g/cm³
- `CalculateQuote()` applies a 20% profit multiplier: `grams × costPerGram × 1.2`
- `POST /orders` re-slices the uploaded file independently (does not reuse the quote price)
- `.3mf` files skip the `--load` config flag; all other extensions load `PRUSA_SLICER_CONFIG`

**Shared types** — repository interfaces and domain structs live in `types/types.go`.

**Config** — loaded from `.env` via `config/env.go` (godotenv). Required vars: `DB_*`, `PORT`, `UPLOAD_PATH`, `PRUSA_SLICER_PATH` (default `/prusa-slicer/AppRun`), `PRUSA_SLICER_CONFIG` (default `/kalleprint/config/prusa-slicer.ini`).

**Migrations** — `cmd/migrate/migrations/` numbered pairs (`N_name.up.sql` / `N_name.down.sql`), run via golang-migrate through the Makefile.

**UUIDs** — stored as binary in MySQL; queries use `UUID_TO_BIN()` / `BIN_TO_UUID()`.

**Error responses** — use `utils.WriteError(w, status, err)` for consistent JSON `{"error": "..."}` responses.

### Client

File-based routing via TanStack Router under `src/routes/`. The main page (`routes/index.tsx`) implements a three-stage state machine:

1. `upload` — user selects file + filament, triggers `POST /quote`
2. `quote` — displays price; user can switch filament (re-quotes inline) or proceed to order form
3. `confirmed` — order placed successfully

Key files:
- `src/lib/api.ts` — Axios instance (`VITE_API_URL` base); typed API functions; error message extracted from `response.data.error`
- `src/lib/queries.ts` — TanStack Query hooks (`useFilaments`, `useGetQuote` mutation, `usePlaceOrder` mutation)
- `src/components/` — page sections (`UploadSection`, `QuoteSection`, `ConfirmationSection`, etc.)
- Styling via Tailwind CSS 4

### Admin

Separate TanStack Start app in `admin/` using pnpm. Routes under `admin/src/routes/`.

### Docker

- **Dev:** `docker-compose.yml` — API uses `Dockerfile.dev` with Air hot reload; client uses `Dockerfile.dev` with Vite HMR; source directories are bind-mounted.
- **Prod:** `docker-compose.prod.yml` — multi-stage build; Go binary compiled statically; tests run during build.
