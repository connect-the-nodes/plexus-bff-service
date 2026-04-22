# Praxis Schema Mapping

This document describes how the newly implemented admin and domain APIs map onto the live PostgreSQL schema in `praxis_schema`.

## Scope

Covered services:

- admin users
- admin groups
- admin permissions
- domains

## Runtime Assumptions

- `DB_SCHEMA` selects the active PostgreSQL schema
- for Praxis deployments, `DB_SCHEMA=praxis_schema`
- API IDs remain strings for compatibility
- PostgreSQL IDs are `bigint` identity columns
- the repository layer converts between API string IDs and database numeric IDs
- when using an existing Praxis schema, `DB_AUTO_MIGRATE=false` should be used

## Configuration Inputs

The current Praxis-backed repository uses:

- `DB_SCHEMA`
- `DB_CUSTOMER_ACCOUNT_ID`
- `DB_DEFAULT_ACTOR_ID`
- `DB_DEFAULT_ACTOR_NAME`

Current usage:

- `DB_SCHEMA` controls schema-qualified table names
- `DB_CUSTOMER_ACCOUNT_ID` is written into Praxis rows such as `admin_user`, `admin_group`, and `domain`
- `DB_DEFAULT_ACTOR_ID` is used as `created_by_user_id` for domain creation when set
- `DB_DEFAULT_ACTOR_NAME` is used as `created_by_display_name`

## Users

API endpoints:

- `GET /api/admin/users`
- `POST /api/admin/users`
- `GET /api/admin/users/{id}`
- `PUT /api/admin/users/{id}`
- `DELETE /api/admin/users/{id}`

Primary tables:

- `praxis_schema.admin_user`
- `praxis_schema.admin_user_group`

Column mapping:

- `PortalUser.ID` <-> `admin_user.admin_user_id`
- `PortalUser.Username` <-> `admin_user.username`
- `PortalUser.PasswordHash` <-> `admin_user.password_value`
- `PortalUser.Email` <-> `admin_user.email_address`
- `PortalUser.DisplayName` <-> `admin_user.display_name`
- `PortalUser.Role` <-> `admin_user.portal_role`
- `PortalUser.Active` <-> `admin_user.is_active`
- `PortalUser.GroupIDs` <-> `admin_user_group.admin_group_id`

Additional Praxis-managed fields:

- `customer_account_id` comes from `DB_CUSTOMER_ACCOUNT_ID`
- `external_id` is generated in the repository
- `created_at` and `updated_at` are managed by PostgreSQL/default SQL expressions

Important note:

- existing sample rows in `admin_user.password_value` are plain text
- new writes from this Go service use bcrypt hashes

## Groups

API endpoints:

- `GET /api/admin/groups`
- `POST /api/admin/groups`
- `GET /api/admin/groups/{id}`
- `PUT /api/admin/groups/{id}`
- `DELETE /api/admin/groups/{id}`

Primary tables:

- `praxis_schema.admin_group`
- `praxis_schema.admin_group_permission`

Column mapping:

- `PortalGroup.ID` <-> `admin_group.admin_group_id`
- `PortalGroup.Name` <-> `admin_group.group_name`
- `PortalGroup.Description` <-> `admin_group.group_description`
- `PortalGroup.PermissionIDs` <-> `admin_group_permission.admin_permission_id`

Additional Praxis-managed fields:

- `customer_account_id` comes from `DB_CUSTOMER_ACCOUNT_ID`
- `external_id` is generated in the repository
- `group_slug` is generated from the group name

## Permissions

API endpoints:

- `GET /api/admin/permissions`
- `POST /api/admin/permissions`
- `GET /api/admin/permissions/{id}`
- `PUT /api/admin/permissions/{id}`
- `DELETE /api/admin/permissions/{id}`

Primary table:

- `praxis_schema.admin_permission`

Column mapping:

- `PermissionDefinition.ID` <-> `admin_permission.admin_permission_id`
- `PermissionDefinition.Name` <-> `admin_permission.permission_name`
- `PermissionDefinition.Description` <-> `admin_permission.permission_description`

Repository-generated Praxis fields:

- `permission_code` is generated from the name
- `permission_category` is currently defaulted to `Custom`

## Domains

API endpoints:

- `GET /api/admin/domains/workspace`
- `GET /api/admin/domains`
- `POST /api/admin/domains`
- `GET /api/admin/domains/approved`
- `GET /api/admin/domains/{id}`
- `PUT /api/admin/domains/{id}`
- `DELETE /api/admin/domains/{id}`
- `POST /api/admin/domains/{id}/decision`

Primary tables:

- `praxis_schema.domain`
- `praxis_schema.domain_review_comment`

Column mapping:

- `RegisteredDomain.ID` <-> `domain.domain_id`
- `RegisteredDomain.Name` <-> `domain.domain_name`
- `RegisteredDomain.Description` <-> `domain.domain_description`
- `RegisteredDomain.OwnerGroupID` <-> `domain.owner_group_id`
- `RegisteredDomain.Status` <-> `domain.lifecycle_status`
- `RegisteredDomain.Review[]` <-> `domain_review_comment`

Repository-generated Praxis fields:

- `customer_account_id` comes from `DB_CUSTOMER_ACCOUNT_ID`
- `external_id` is generated in the repository
- `domain_code` is generated from the domain name
- `owner_role_name` is generated as `<domain name> Approver`
- `created_by_user_id` comes from `DB_DEFAULT_ACTOR_ID` when set
- `created_by_display_name` comes from `DB_DEFAULT_ACTOR_NAME`

Status mapping:

- API `mode=DRAFT` -> `DRAFT`
- API `mode=SUBMIT` -> `PENDING_APPROVAL`
- review `APPROVED` -> `APPROVED`
- review `REJECTED` -> `REJECTED`

## Domain Workspace

API endpoint:

- `GET /api/admin/domains/workspace`

Current implementation:

- loads owner groups from `admin_group`
- loads domains from `domain`
- returns a combined DTO for the UI workflow

## Current Gaps

- `domain.metadata` is accepted by the API but not persisted because the live Praxis `domain` table has no metadata column
- `permission_code` and `permission_category` are not caller-controlled in the current DTO
- `domain_code`, `owner_role_name`, and `external_id` are generated server-side
- actor attribution currently comes from config defaults, not authenticated user context
