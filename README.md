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

## Notes

- Verified locally on this machine with `go1.26.2`:
  - `go test ./...`
  - `go build ./cmd/plexus-bff-service`
  - runtime smoke test on `http://localhost:8080/actuator/health`
- See [docs/parity-gap-report.md](../plexus-bff-service-go/docs/parity-gap-report.md) for remaining parity caveats and follow-up items.
