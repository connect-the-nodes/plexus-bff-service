# Handoff

## Current Status

The project has been extended with PostgreSQL-backed admin services for:

- users
- groups
- permissions
- domains

The implementation is in place and `go test ./...` passes.

## Completed Work

- Added PostgreSQL configuration support in `configs/` and `internal/app/config/config.go`
- Added PostgreSQL connection bootstrap and embedded migration runner
- Added SQL migration for:
  - `users`
  - `user_groups`
  - `groups`
  - `group_permissions`
  - `permissions`
  - `domains`
  - `domain_reviews`
  - `schema_migrations`
- Added DTOs, models, mappers, repository contracts, service implementations, and controllers
- Added route wiring for:
  - `GET/POST /api/admin/users`
  - `GET/PUT/DELETE /api/admin/users/{id}`
  - `GET/POST /api/admin/groups`
  - `GET/PUT/DELETE /api/admin/groups/{id}`
  - `GET/POST /api/admin/permissions`
  - `GET/PUT/DELETE /api/admin/permissions/{id}`
  - `GET /api/admin/domains/workspace`
  - `GET/POST /api/admin/domains`
  - `GET /api/admin/domains/approved`
  - `GET/PUT/DELETE /api/admin/domains/{id}`
  - `POST /api/admin/domains/{id}/decision`
- Added in-memory fallback repositories for tests and non-DB runs
- Added documentation in `docs/postgresql-admin-services.md`
- Updated `README.md`

## Important Business Rules Implemented

- Passwords are stored as bcrypt hashes
- Users can belong to many groups
- Groups can have many permissions
- Each domain has one owner group
- Domain `mode=DRAFT` maps to `DRAFT`
- Domain `mode=SUBMIT` maps to `PENDING_APPROVAL`
- Domain review decisions map to `APPROVED` or `REJECTED`
- Only draft domains can be deleted
- Only pending-review domains can be reviewed

## Verified Against Local Praxis Database

The admin/domain services were verified against the user's live local PostgreSQL instance with:

- host: `localhost`
- port: `5432`
- database: `praxis_db`
- schema: `praxis_schema`
- app user: `praxis_user`

Verified successfully:

- database connection with real local credentials
- live reads for users/groups/permissions/domains
- live create/delete cycle for permissions
- live create/delete cycle for groups
- live create/delete cycle for users
- live create/delete cycle for domains
- counts returning to baseline after cleanup

Observed baseline after verification:

- users: `6`
- groups: `6`
- permissions: `9`
- domains: `2`

Still not verified in this workspace:

- `DB_AUTO_MIGRATE=true` execution against a `public` schema deployment
- end-to-end domain review workflow against a pending-approval row in Praxis data
- authenticated actor propagation from JWT context into created-by/review attribution

## Known Blocker

`gofmt` could not be executed on this machine because local Application Control blocked `gofmt.exe`.

This is an environment restriction, not a compile/test failure.

## Next Recommended Step

Resume by hardening the Praxis integration:

1. Decide whether `admin_user.password_value` should be normalized to bcrypt for existing rows
2. Decide whether the API contract should expose Praxis-specific fields like `permission_code`, `permission_category`, `domain_code`, and `owner_role_name`
3. Add repository-level tests against a temporary Praxis-compatible schema
4. Wire actor attribution from authenticated user context instead of config defaults
5. Verify the domain review flow using a real `PENDING_APPROVAL` domain
6. Replace manual `static/openapi.yaml` maintenance with a repeatable OpenAPI generation or synchronization workflow so Swagger UI stays aligned automatically

## Key Files

- `docs/postgresql-admin-services.md`
- `docs/praxis-schema-mapping.md`
- `docs/runtime-verification.md`
- `internal/app/database/postgres.go`
- `internal/app/database/migrator.go`
- `internal/app/database/migrations/001_admin_identity.sql`
- `internal/app/app/app.go`
- `internal/app/httpx/router.go`
- `internal/app/repository/impl/admin_postgres_repository.go`
- `internal/app/service/impl/admin_user_service.go`
- `internal/app/service/impl/admin_group_service.go`
- `internal/app/service/impl/admin_permission_service.go`
- `internal/app/service/impl/admin_domain_service.go`

## Suggested Prompt For Next Session

Use:

`Continue from docs/handoff.md and docs/praxis-schema-mapping.md. Harden the Praxis PostgreSQL integration and close the remaining mapping and audit gaps.`
