# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Backend for the Protacc website — a Go Fiber REST API over PostgreSQL (Neon),
deployed on Render. Covers user auth, a service catalog, cart/orders, Razorpay
payments, lead/contact capture, admin analytics, and configurable system settings.

## Commands

```bash
go run main.go            # Run the server (defaults to :8080, override with PORT)
go build -o protacc-backend   # Build the binary (Procfile runs ./protacc-backend)
go mod tidy               # Sync dependencies
go vet ./...              # Vet

# Migrations (golang-migrate CLI required: brew install golang-migrate)
./scripts/migrate.sh up                  # Apply migrations; reads DATABASE_URL from .env
./scripts/migrate.sh down                # Roll back
./scripts/migrate.sh create <name>       # New sequential .up.sql/.down.sql pair in db/migrations
./scripts/migrate.sh force <version>     # Force a dirty version
./scripts/migrate.sh up "<database_url>" # 2nd arg overrides DATABASE_URL (e.g. a prod branch)
```

There are no tests in this repo. `curl localhost:8080/ping` should return `pong`.

## Architecture

Strict four-layer flow, one file per domain in each layer. A request moves:
`handler/ → service/ → repository/ → db.Pool`. Domains: user, service, cart,
order, order_document, payment, settings, analytics, lead, contact (plus
`mail_service` which has no repo).

- **`handler/`** — Fiber HTTP handlers. Parse/validate input (`go-playground/validator`),
  read auth context via `c.Locals("userId"|"email"|"role")`, return JSON. No business logic.
- **`service/`** — business logic; orchestrates one or more repositories.
- **`repository/`** — raw SQL against `pgx`. Each holds a `*pgxpool.Pool`. Multi-step
  writes (e.g. order + order items) use `tx := db.Begin(...)` with `defer tx.Rollback`.
- **`models/`** — request/response and DB structs per domain.

**Dependency wiring is centralized and manual** in `app/factories.go`
(`NewFactories`): it builds repos → services → handlers in order and returns a
`Factories` struct. Add a new domain by constructing its repo/service/handler
there, adding the handler to `Factories`, and registering routes.

**Routes** live entirely in `app/routes.go` under `/api/v1`. Three access tiers:
- Public (no middleware): auth, service browsing, public settings, lead/contact creation, payment webhook.
- `middleware.Protected()`: validates `Bearer` JWT, populates `c.Locals`.
- `middleware.Protected()` + `middleware.RequireRole("admin")`: admin group.

**DB connection** is a single global `db.Pool` (`db/db.go`), initialized in
`NewFactories`. Note: prepared statements are disabled
(`QueryExecModeExec`) to avoid conflicts with Neon's pooler — keep this when
touching pool config.

## Conventions & gotchas

- **Module path is `github.com/aman4411/protacc-backend`** — internal imports use it.
- **Env var is `DATABASE_URL`**, not `DB_URL` (the README is stale on this). Required
  vars: `DATABASE_URL`, `JWT_SECRET`, `RESEND_API_KEY`, `FROM_EMAIL`, `FRONTEND_URL`,
  `RAZORPAY_KEY_ID`/`_KEY_SECRET`/`_WEBHOOK_SECRET`, optional `PORT`, `ENVIRONMENT`.
- **Roles:** middleware checks the string `"admin"`. `common/constants.go` defines
  `RoleCustomer`/`RoleStaff`/`RoleAdmin` but `RoleAdmin = "admin"` is the one enforced.
- **Payments** use Razorpay; if credentials are absent the `PaymentService` logs a
  warning and runs with a nil client rather than failing startup. Webhook signatures
  are HMAC-SHA256 verified.
- **Email** goes through Resend's HTTP API (`mail_service.go`), not SMTP.
- **Order numbers** are generated as `ORD<yyyymmddhhmmss>`; order status is an enum
  (see the `update_order_status_enum` migrations) tracking a booking→completion flow.
- **Google Drive docs**: order documents store shared Drive links, parsed/normalized
  by `utils/ParseGoogleDriveURL` in `utils/google_drive.go`.
- **CORS** allowlist is hardcoded in `routes.go` (protacc.in, the Netlify domain,
  localhost:3000) — update there when origins change.
