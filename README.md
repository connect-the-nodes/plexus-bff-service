# plexus-bff-service-go

Go migration of `plexus-bff-service` with package structure and runtime behavior kept intentionally close to the Java service for side-by-side comparison.

## Scope

This repository mirrors:

- API routes and response envelope conventions
- profile-based configuration
- feature flag loading and optional session caching
- observability proxy endpoints
- auth redirect endpoints
- correlation ID propagation
- optional Redis-backed session behavior
- container and build artifacts
- PostgreSQL-backed admin user/group/permission/domain services

## Layout

- `cmd/plexus-bff-service`: application entrypoint
- `internal/app/config`: configuration loading and profile overlays
- `internal/app/controller`: HTTP handlers aligned with Java controllers
- `internal/app/service`: service layer
- `internal/app/repository`: repository layer
- `internal/app/model`, `dto`, `mapper`, `feature`, `security`, `tracing`: Java-aligned support packages
- `configs`: migrated application profile YAML files
- `docs`: migration notes, parity map, and known gaps

## Build

```powershell
go test ./...
go build ./cmd/plexus-bff-service
```

## Run

```powershell
$env:SPRING_PROFILES_ACTIVE="local"
go run ./cmd/plexus-bff-service
```

Default local URLs:

- `http://localhost:8080/actuator/health`
- `http://localhost:8080/swagger-ui.html`
- `http://localhost:8080/v3/api-docs`
- `http://localhost:8080/api/admin/users`
- `http://localhost:8080/api/admin/groups`
- `http://localhost:8080/api/admin/permissions`
- `http://localhost:8080/api/admin/domains`

## PostgreSQL

The admin user/group/permission/domain services support PostgreSQL via embedded migrations.

Typical local setup:

```powershell
$env:SPRING_PROFILES_ACTIVE="local"
$env:DB_INTEGRATION_ENABLED="true"
$env:DB_AUTO_MIGRATE="false"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_NAME="praxis_db"
$env:DB_USER="praxis_user"
$env:DB_PASSWORD="praxis_user"
$env:DB_SCHEMA="praxis_schema"
$env:DB_CUSTOMER_ACCOUNT_ID="1"
$env:DB_DEFAULT_ACTOR_ID="1"
$env:DB_DEFAULT_ACTOR_NAME="Portal Administrator"
go run ./cmd/plexus-bff-service
```

On startup the application will:

- open a PostgreSQL connection pool
- run embedded migrations when `DB_AUTO_MIGRATE=true` and `DB_SCHEMA=public`
- expose the admin APIs backed by PostgreSQL, including existing Praxis schemas such as `praxis_schema`

See [docs/postgresql-admin-services.md](docs/postgresql-admin-services.md) for the implementation details.
See [docs/praxis-schema-mapping.md](docs/praxis-schema-mapping.md) for the live Praxis table mapping.
See [docs/runtime-verification.md](docs/runtime-verification.md) for the runbook and verification flow.

## Swagger UI

The local Swagger UI is available at:

- `http://localhost:8080/swagger-ui.html`

The OpenAPI source document is served at:

- `http://localhost:8080/v3/api-docs`

Maintenance rule:

- when adding or changing a locally testable endpoint, update `internal/app/httpx/router.go` and `static/openapi.yaml` in the same change so Swagger UI stays aligned with runtime behavior

## Notes

- Verified locally on this machine with `go1.26.2`:
  - `go test ./...`
  - `go build ./cmd/plexus-bff-service`
  - runtime smoke test on `http://localhost:8080/actuator/health`
- See [docs/parity-gap-report.md](../plexus-bff-service-go/docs/parity-gap-report.md) for remaining parity caveats and follow-up items.
