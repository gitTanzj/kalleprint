# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build & run
make build          # compile to bin/server
make run            # build + run

# Tests
make test           # go test -v ./...

# Database migrations
make migrate-up     # apply pending migrations
make migrate-down   # rollback migrations
make migration NAME=<name>  # create new migration pair

# Development (Docker with hot reload)
docker compose up
```

## Architecture

Go REST API for a 3D printing service. Entry point is `cmd/main.go`, which initializes the DB connection and HTTP server on port 8080.

**Service pattern:** Each service lives under `services/<name>/` and contains:
- `routes.go` — registers HTTP handlers on the router
- `repository.go` — database queries

Services are registered in `cmd/api/api.go` by calling `NewRouter()` on each service and passing the `*mux.Router`.

**Current services:** `health`, `filaments`, `orders`, `prints`, `quotes` (file upload), `jobs` (WIP — routes defined, handlers empty).

**Types** shared across services live in `types/types.go`.

**Config** is loaded from `.env` via `config/env.go`; see `.env.example` for required variables (DB credentials, PORT, UPLOAD_PATH).

**Migrations** live in `cmd/migrate/migrations/` numbered sequentially (e.g. `1_customers-table.up.sql`). Run via `golang-migrate`; the migrate binary is invoked by the Makefile.

**Docker:** `docker-compose.yml` runs MySQL + the API with `air` hot reload. `Dockerfile.prod` / `docker-compose.prod.yml` handle production multi-stage builds.
