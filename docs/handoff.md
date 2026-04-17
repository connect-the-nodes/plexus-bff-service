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
- Domain `mode=SUBMIT` maps to `PENDING_REVIEW`
- Domain review decisions map to `APPROVED` or `REJECTED`
- Only draft domains can be deleted
- Only pending-review domains can be reviewed

## Not Yet Verified

These still need to be verified against the user's actual local PostgreSQL instance:

- database connection with real local credentials
- auto-migration execution on startup
- live CRUD behavior for users/groups/permissions/domains
- domain review flow against PostgreSQL
- referential integrity behavior in real DB execution

## Known Blocker

`gofmt` could not be executed on this machine because local Application Control blocked `gofmt.exe`.

This is an environment restriction, not a compile/test failure.

## Next Recommended Step

Resume by performing real PostgreSQL verification:

1. Confirm PostgreSQL is running locally
2. Create the target database if needed, for example `plexus_bff`
3. Run the service with:
   - `SPRING_PROFILES_ACTIVE=local`
   - `DB_INTEGRATION_ENABLED=true`
   - `DB_HOST`
   - `DB_PORT`
   - `DB_NAME`
   - `DB_USER`
   - `DB_PASSWORD`
4. Let migrations run
5. Exercise the new admin endpoints
6. Fix any runtime issues found from real DB interaction
7. Optionally add sample seed data and API examples

## Key Files

- `docs/postgresql-admin-services.md`
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

`Continue from docs/handoff.md. Verify the PostgreSQL-backed admin user/group/permission/domain services against my local Postgres and finish the remaining runtime validation.`
