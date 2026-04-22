# PostgreSQL Admin Services

This document covers the PostgreSQL-backed implementation added for:

- admin users
- admin groups
- admin permissions
- domains

For Praxis-specific repository mapping and a live verification runbook, also see [praxis-schema-mapping.md](praxis-schema-mapping.md) and [runtime-verification.md](runtime-verification.md).

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

These endpoints are exposed for browser-based local testing through:

- Swagger UI: `/swagger-ui.html`
- OpenAPI document: `/v3/api-docs`

## Runtime Model

The application supports two execution paths for these services:

- PostgreSQL-backed repositories when `DB_INTEGRATION_ENABLED=true`
- in-memory repositories when `DB_INTEGRATION_ENABLED=false`

The in-memory path keeps the service runnable for tests and local development. The PostgreSQL path is the intended production path.

For existing Praxis databases, the expected setup is:

- `DB_SCHEMA=praxis_schema`
- `DB_AUTO_MIGRATE=false`

## Configuration

Environment variables:

- `DB_INTEGRATION_ENABLED`
- `DB_AUTO_MIGRATE`
- `DB_SCHEMA`
- `DB_CUSTOMER_ACCOUNT_ID`
- `DB_DEFAULT_ACTOR_ID`
- `DB_DEFAULT_ACTOR_NAME`
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

When wiring to an existing enterprise schema such as `praxis_schema`, set `DB_AUTO_MIGRATE=false` and `DB_SCHEMA=praxis_schema`. The Go repositories will then operate against the existing table names in that schema instead of the embedded migration tables.

Key relationships:

- users to groups: many-to-many
- groups to permissions: many-to-many
- domains to groups: many-to-one through `owner_group_id`
- domains to reviews: one-to-many

## Domain Status Rules

Domain lifecycle mapping:

- `mode=DRAFT` -> `DRAFT`
- `mode=SUBMIT` -> `PENDING_APPROVAL`
- review decision `APPROVED` -> `APPROVED`
- review decision `REJECTED` -> `REJECTED`

Business constraints:

- only draft domains can be deleted
- only pending-review domains can be reviewed
- user `groupIds` must reference existing groups
- group `permissionIds` must reference existing permissions
- domain `ownerGroupId` must reference an existing group

## Praxis-Specific Mapping Notes

When `DB_SCHEMA=praxis_schema`, the repository layer maps the API onto:

- `admin_user`
- `admin_group`
- `admin_permission`
- `admin_user_group`
- `admin_group_permission`
- `domain`
- `domain_review_comment`

The API still exposes string IDs, while the Praxis database uses `bigint` identity columns. The repository converts between those representations.

The current implementation also generates several Praxis fields server-side:

- `external_id`
- `group_slug`
- `permission_code`
- `domain_code`
- `owner_role_name`

## Security Notes

- passwords are stored as bcrypt hashes
- password hashes are never returned in API responses
- the DTO still includes a `password` field for request compatibility, but responses intentionally omit it

Important Praxis caveat:

- existing sample rows in `praxis_schema.admin_user.password_value` may still be plain text
- newly created or updated users through this Go service are written as bcrypt hashes

## Current Limitations

- `domain.metadata` is accepted by the API but not persisted to PostgreSQL because the live Praxis `domain` table has no metadata column
- `permission_code` and `permission_category` are not caller-controlled in the current DTO and are generated/defaulted in the repository
- `domain_code`, `owner_role_name`, and `external_id` are generated server-side rather than being exposed through the API contract
- actor attribution for domain creation currently comes from configuration defaults, not the authenticated user context

## Migrations

Migrations are embedded under `internal/app/database/migrations`.

When `DB_AUTO_MIGRATE=true`, migrations run automatically during application startup before routes are served.
