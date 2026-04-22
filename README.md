# EduCRM Backend

REST API for EduCRM: **Go 1.26**, **Gin**, **GORM**, **PostgreSQL**, **JWT** access + opaque refresh tokens, **Swagger**, structured logging, and Docker. Business logic lives behind **use cases** and **repository interfaces** (clean architecture).

## Prerequisites

- **Go** 1.26+
- **PostgreSQL** 16+ (local, managed, or Docker Compose)
- **Docker** (optional, for Compose and image builds)

## Quick setup

1. Clone the repository and `cd` into it.
2. Copy environment template: `cp .env.example .env` and edit values (or export variables in your shell). The API does **not** load `.env` by itself; use your shell, systemd, Kubernetes secrets, or Compose `environment:` blocks.
3. Create a database matching `DB_*`.
4. Apply schema (choose one path):
   - **Development:** leave `AUTO_MIGRATE=true` (default) and start the API — GORM syncs models.
   - **Production-style:** set `AUTO_MIGRATE=false`, run SQL migrations (`make migrate-up`), then start the API.
5. Create the first **super_admin** (once): `make seed` with `SEED_SUPER_ADMIN_EMAIL` and `SEED_SUPER_ADMIN_PASSWORD` set (see [Seed](#seed-initial-super_admin)).
6. Open Swagger: `http://localhost:8080/swagger/index.html` (when `ENABLE_SWAGGER` is on).

## Environment separation

| `APP_ENV`     | Gin mode   | Notes |
|---------------|------------|--------|
| `development` | debug      | Relaxed defaults; CORS may use `*` when `CORS_ALLOWED_ORIGINS` is unset. |
| `staging`     | debug      | `ValidateForAPI` enforces strong `JWT_SECRET` (32+ chars, no placeholders). |
| `production`  | release    | Strong JWT, `DB_SSLMODE` not `disable`, explicit `CORS_ALLOWED_ORIGINS` (no `*`), Swagger off unless `ENABLE_SWAGGER=true`. |

Configuration is loaded only from **environment variables** (`internal/config`). See `.env.example` for names and defaults.

### Example environment (abbreviated)

```bash
APP_ENV=development
LOG_LEVEL=info
HTTP_PORT=8080
SHUTDOWN_TIMEOUT=30s
AUTO_MIGRATE=true

DB_HOST=localhost
DB_PORT=5432
DB_USER=educrm
DB_PASSWORD=educrm
DB_NAME=educrm
DB_SSLMODE=disable

JWT_SECRET=replace-with-at-least-32-random-characters
JWT_ACCESS_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=168h

# Optional: comma-separated origins; credentials require non-wildcard origins
# CORS_ALLOWED_ORIGINS=https://app.example.com
# CORS_ALLOW_CREDENTIALS=false

# Optional: in-process rate limit (per IP)
# RATE_LIMIT_ENABLED=true
# RATE_LIMIT_RPS=100
# RATE_LIMIT_BURST=200
```

Full list: **`.env.example`**.

## Commands

| Command | Description |
|---------|-------------|
| `make run` | Run API (`go run ./cmd/api`). |
| `make build` | Build `bin/api`. |
| `make build-tools` | Build `bin/api`, `bin/migrate`, `bin/seed`. |
| `make test` | `go test ./...` |
| `make vet` | `go vet ./...` |
| `make tidy` | `go mod tidy` |
| `make swag` | Regenerate `docs/` from handler comments. |
| `make migrate-up` | Apply SQL migrations (`cmd/migrate up`). |
| `make migrate-down` | Roll back one migration step. |
| `make migrate-version` | Print current migration version. |
| `make seed` | Create super_admin if email is free (`cmd/seed`). |
| `make docker-up` / `make docker-down` | Compose stack (Postgres + API). |

### Migration run command

From repo root (requires `DB_*` and optional `MIGRATIONS_PATH`, default `migrations`):

```bash
go run ./cmd/migrate up
go run ./cmd/migrate down
go run ./cmd/migrate version
```

In Docker (override entrypoint; same `DB_*` as the API service):

```bash
docker compose run --rm --entrypoint /app/migrate api up
```

### Seed (initial super_admin)

Required environment variables:

- `SEED_SUPER_ADMIN_EMAIL` — must look like an email (used for login).
- `SEED_SUPER_ADMIN_PASSWORD` — stored with bcrypt.

Optional: `SEED_SUPER_ADMIN_PHONE`.

```bash
export SEED_SUPER_ADMIN_EMAIL='admin@school.edu'
export SEED_SUPER_ADMIN_PASSWORD='your-secure-password'
make seed
```

If that email already exists, the command exits successfully without changes.

Docker:

```bash
docker compose run --rm --entrypoint /app/seed \
  -e SEED_SUPER_ADMIN_EMAIL=admin@school.edu \
  -e SEED_SUPER_ADMIN_PASSWORD='your-secure-password' \
  api
```

## Project structure

```
cmd/
  api/          # HTTP server entrypoint
  migrate/      # golang-migrate CLI wrapper (SQL in ./migrations)
  seed/         # One-shot super_admin bootstrap
internal/
  app/          # Composition root, lifecycle, graceful shutdown
  config/       # Env config, validation (production/staging), CORS helpers
  database/     # Postgres + GORM pool
  delivery/http/# Gin router, handler adapters, DTOs
  middleware/   # Recovery, request ID, logging, CORS, rate limit, JWT, RBAC
  domain/       # Entities and invariants
  usecase/      # Application services
  repository/   # Ports; postgres/ implements GORM
  rbac/         # Permission matrix
  ai/, notify/, storage/ ...
migrations/     # Versioned .up.sql / .down.sql
docs/           # Generated Swagger
pkg/            # jwt, logger, response (shared, low-level)
```

Dependency direction: **HTTP → use case → domain**; repositories are injected in `internal/app`.

## Run locally

```bash
go mod download
go run ./cmd/api
```

- **Liveness** (process up, no DB): `GET /health`
- **Readiness** (DB ping): `GET /api/v1/health` — returns **503** if the database is down.

## Run with Docker Compose

```bash
make docker-up
```

API: `http://localhost:8080`. Postgres is exposed on `5432` with user/password/db `educrm`.

Image includes **`/app/api`**, **`/app/migrate`**, **`/app/seed`**, and **`/app/migrations`**. Non-root user, **HEALTHCHECK** on `/health`, minimal runtime with `curl`.

## Production-oriented behavior

- **Secure config:** `config.ValidateForAPI()` runs in `cmd/api` for staging/production rules (JWT, TLS mode, CORS).
- **Graceful shutdown:** SIGINT/SIGTERM stops new traffic; in-flight requests get up to `SHUTDOWN_TIMEOUT` to finish (`internal/app`).
- **CORS:** Configurable origins and optional credentials (`internal/middleware/cors.go`).
- **Rate limiting:** Optional per-IP token bucket (`RATE_LIMIT_*`); safe for single instance — use a shared store for horizontal scale.
- **Request logging:** JSON access lines with `request_id`, `client_ip`, latency; paths in `LOG_HTTP_SKIP_PATHS` are skipped (default health routes).
- **Panic recovery:** Logs stack trace and returns generic 500 JSON (`internal/middleware/recovery.go`).

## Tests

```bash
make test
# or
go test ./...
```

Handler, middleware, RBAC, domain, and selected use-case tests use table-driven patterns where practical.

## API surface

Interactive docs: **`/swagger/index.html`** when Swagger is enabled.

**Frontend / Cursor:** see **`docs/FRONTEND_API.md`** (envelope, auth, base URL, endpoint map, Cursor rule snippet).

Domains include authentication, users, teachers, rooms, groups, schedules, attendance, grades, payments, files, notifications, dashboard, and AI analytics. Roles: `super_admin`, `admin`, `teacher`, `student` (see Swagger for paths and RBAC).

## Module path

Go module: `github.com/educrm/educrm-backend`. Forks can `go mod edit -module <path>` and fix imports.

## License

See repository license file if present.
# EduCRM
