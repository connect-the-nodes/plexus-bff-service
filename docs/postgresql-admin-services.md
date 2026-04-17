# PostgreSQL Admin Services

This document covers the PostgreSQL-backed implementation added for:

- admin users
- admin groups
- admin permissions
- domains

## Endpoints

Implemented endpoints:

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

## Runtime Model

The application supports two execution paths for these services:

- PostgreSQL-backed repositories when `DB_INTEGRATION_ENABLED=true`
- in-memory repositories when `DB_INTEGRATION_ENABLED=false`

The in-memory path keeps the service runnable for tests and local development. The PostgreSQL path is the intended production path.

## Configuration

Environment variables:

- `DB_INTEGRATION_ENABLED`
- `DB_AUTO_MIGRATE`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`
- `DB_SSLMODE`
- `DB_MAX_CONNS`
- `DB_MIN_CONNS`
- `DB_MAX_CONN_LIFETIME_MINUTES`

Default local values are defined in `configs/application.yml`, with `application-local.yml` enabling DB integration by default.

## Database Schema

Embedded migrations create:

- `users`
- `user_groups`
- `groups`
- `group_permissions`
- `permissions`
- `domains`
- `domain_reviews`
- `schema_migrations`

Key relationships:

- users to groups: many-to-many
- groups to permissions: many-to-many
- domains to groups: many-to-one through `owner_group_id`
- domains to reviews: one-to-many

## Domain Status Rules

Domain lifecycle mapping:

- `mode=DRAFT` -> `DRAFT`
- `mode=SUBMIT` -> `PENDING_REVIEW`
- review decision `APPROVED` -> `APPROVED`
- review decision `REJECTED` -> `REJECTED`

Business constraints:

- only draft domains can be deleted
- only pending-review domains can be reviewed
- user `groupIds` must reference existing groups
- group `permissionIds` must reference existing permissions
- domain `ownerGroupId` must reference an existing group

## Security Notes

- passwords are stored as bcrypt hashes
- password hashes are never returned in API responses
- the DTO still includes a `password` field for request compatibility, but responses intentionally omit it

## Migrations

Migrations are embedded under `internal/app/database/migrations`.

When `DB_AUTO_MIGRATE=true`, migrations run automatically during application startup before routes are served.
