# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

This is the **admin panel** for the Kalleprint 3D printing service. See the parent `../CLAUDE.md` for the full-stack overview (Go server, customer client, MySQL, PrusaSlicer).

## Commands

Package manager is **npm** (despite the parent CLAUDE.md mentioning pnpm — that is stale; `package-lock.json` and `Dockerfile.dev` use npm).

```bash
npm run dev       # Vite dev server on :3000
npm run build     # production build → dist/
npm test          # vitest run
npm run lint      # eslint
npm run format    # prettier --write + eslint --fix
npm run check     # prettier --check (no writes)
```

Run a single test file: `npx vitest run path/to/file.test.ts`.

## Architecture

**Framework:** TanStack Start (SSR-capable React) + TanStack Router (file-based routing) + TanStack Query. React 19, Tailwind 4 (via `@tailwindcss/vite`), Vite 8.

**Route tree** is auto-generated into `src/routeTree.gen.ts` by `@tanstack/router-plugin` — do not edit it by hand; add/move files under `src/routes/` and the plugin regenerates on dev/build.

**Path aliases:** `#/*` and `@/*` both resolve to `src/*` (defined in `tsconfig.json` and `package.json#imports`). Existing code uses relative imports, but either alias is valid.

### Auth flow

JWT auth, token stored in `localStorage` under `admin_token` (+ `admin_token_expires_at`). All API requests automatically attach `Authorization: Bearer <token>` via an axios interceptor in `src/lib/api.ts`.

The route guard lives in `src/routes/__root.tsx#beforeLoad`: any non-`/login` route redirects to `/login` if `isTokenValid()` is false, and `/login` redirects to `/` if already authed. When adding new routes, they are auth-gated by default through this root guard — no per-route work needed.

The `Sidebar` is rendered for every non-login route by `RootLayout` in `__root.tsx`; the login page renders bare.

### API layer

`src/lib/api.ts` — axios instance pointed at `${VITE_API_URL ?? 'http://localhost:8080'}/api/v1/admin`. Note the `/admin` suffix: admin endpoints are a separate API surface from the customer-facing `/api/v1` routes used by the `client/` app.

`src/lib/queries.ts` — TanStack Query hooks wrapping the api functions. Mutations invalidate by query-key root (e.g. `['orders']`, `['filaments']`) on success.

**Quirk:** `fetchOrder(id)` calls `GET /orders` and filters client-side rather than hitting a `/orders/:id` endpoint — the backend does not yet expose a single-order GET. If the orders list grows large, this becomes a problem.

### Order status state machine

The valid statuses (used in both `routes/orders.tsx` filter buttons and `routes/orders/$orderId.tsx` select) are:

`pending` → `printing` → `waiting-for-shipment` → `shipped` → `done`

Keep these two lists in sync with the backend's allowed values when changing them.

### Styling convention

Hard-edged "neobrutalist" look: `border-2 border-black`, flat white/gray backgrounds, no rounded corners, no shadows. When adding UI, match this — don't introduce `rounded-*`, `shadow-*`, or color gradients unless intentionally redesigning.
