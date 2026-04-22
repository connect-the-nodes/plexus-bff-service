# Runtime Verification

This runbook describes how to verify the PostgreSQL-backed admin and domain services against a live Praxis database.

## Verified Environment

Verified locally in this workspace on `2026-04-21` against:

- host: `localhost`
- port: `5432`
- database: `praxis_db`
- schema: `praxis_schema`
- app user: `praxis_user`

Baseline rows observed before and after verification:

- `admin_user`: `6`
- `admin_group`: `6`
- `admin_permission`: `9`
- `domain`: `2`

## Required Environment Variables

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
```

## Start The Service

```powershell
go build ./cmd/plexus-bff-service
.\plexus-bff-service.exe
```

Expected behavior:

- service starts on `http://localhost:8080`
- PostgreSQL connection succeeds
- embedded migrations are skipped because `DB_AUTO_MIGRATE=false`

## Read Verification

```powershell
Invoke-RestMethod http://localhost:8080/api/admin/users
Invoke-RestMethod http://localhost:8080/api/admin/groups
Invoke-RestMethod http://localhost:8080/api/admin/permissions
Invoke-RestMethod http://localhost:8080/api/admin/domains
Invoke-RestMethod http://localhost:8080/api/admin/domains/workspace
Invoke-RestMethod http://localhost:8080/api/admin/domains/approved
```

Expected counts from the verified sample data:

- users: `6`
- groups: `6`
- permissions: `9`
- domains: `2`

## Write Verification

Use a temporary suffix and create then delete rows in this order:

1. Create permission
2. Create group referencing that permission
3. Create user referencing that group
4. Create draft domain referencing that group
5. Delete domain
6. Delete user
7. Delete group
8. Delete permission

Observed result in this workspace:

- permission create/delete worked
- group create/delete worked
- user create/delete worked
- domain create/delete worked
- final counts returned to baseline

## Cleanup Verification SQL

```sql
SELECT count(*) AS temp_permissions
FROM praxis_schema.admin_permission
WHERE permission_name LIKE 'Temp Permission %';

SELECT count(*) AS temp_groups
FROM praxis_schema.admin_group
WHERE group_name LIKE 'Temp Group %';

SELECT count(*) AS temp_users
FROM praxis_schema.admin_user
WHERE username LIKE 'temp.user.%';

SELECT count(*) AS temp_domains
FROM praxis_schema.domain
WHERE domain_name LIKE 'Temp Domain %';
```

Expected result:

- all counts are `0`

## Known Caveats

- existing sample values in `praxis_schema.admin_user.password_value` appear to be plain text
- newly created users from this Go service are stored using bcrypt hashes
- `domain.metadata` is accepted by the API but not persisted because the live `domain` table has no metadata column
- some Praxis-specific fields are generated server-side during persistence
