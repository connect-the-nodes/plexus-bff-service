# Integration Marketplace Design

This document describes the control-plane design for a multi-tenant integration marketplace
in Zynchub Hub Service. The goal is to allow many customers (tenants) to subscribe to a
global catalog of connectors and their APIs, without creating bespoke infrastructure per tenant.

Scope
- Global connector catalog (shared by all tenants)
- 1:1 mapping between connector operations and API endpoints
- Tenants cannot override global schemas; they can only select fields and define mappings

---

## Core concepts

### Control plane (configuration + entitlements)
- Who is the tenant?
- Which connectors are enabled?
- Which operations are enabled?
- Which credentials/config are stored for the tenant?
- Which fields are selected and how are they mapped?

### Data plane (execution)
- Executes the integration call based on the control-plane configuration.
- Enforces subscriptions, applies mappings, uses credentials, and emits metrics.

---

## Data plane (runtime execution)

Purpose
- Execute subscribed operations with tenant-specific configuration and mappings.

Runtime flow (sync)
1) Client calls `POST /integrations/{connectorId}/{operationId}`.
2) Hub resolves tenantId from JWT.
3) Hub validates connector + operation subscription = ACTIVE.
4) Hub loads:
   - operation schema
   - tenant mapping + selected fields
   - connector credentials (Secrets Manager reference)
5) Hub routes to connector runner with normalized payload.
6) Runner calls target API and returns result.
7) Hub records metrics + audit, returns response to client.

Runtime flow (async)
1) Hub validates subscription and enqueues job to SQS.
2) Worker consumes job, executes connector, writes result to callback or status store.
3) Hub (or UI) can fetch job status via `GET /jobs/{jobId}`.

Observability
- Emit metrics by tenantId, connectorId, operationId, status, latency.
- Store aggregates for dashboards and alerts.

## API design (Hub Service)

All APIs below are tenant-scoped and authenticated. Tenant identity comes from the
user session/JWT, and is not passed as a path parameter.

### 1) Connector catalog

1. `GET /connectors`
   Purpose:
   - List the global connector catalog
   Returns:
   - connectorId, name, version, status, supportedAuthTypes

2. `GET /connectors/{connectorId}`
   Purpose:
   - Get detailed connector metadata
   Returns:
   - description, categories/tags, supportedAuthTypes
   - link to operations

3. `GET /connectors/{connectorId}/operations`
   Purpose:
   - List operations (APIs) exposed by a connector
   Returns:
   - operationId, method, path, summary, tags

4. `GET /connectors/{connectorId}/operations/{operationId}`
   Purpose:
   - Get full API spec for a connector operation
   Returns:
   - HTTP method, path
   - request schema (fields, types, required, enums)
   - response schema (fields, types, required, enums)
   - examples

### 2) Subscriptions (tenant entitlements)

5. `GET /subscriptions`
   Purpose:
   - List the tenant’s subscribed connectors and operations
   Returns:
   - connector subscription status
   - operation subscription status

6. `POST /subscriptions`
   Purpose:
   - Subscribe tenant to a connector
   Body:
   - connectorId
   Effects:
   - Creates a subscription with status ACTIVE

7. `POST /subscriptions/{connectorId}/operations`
   Purpose:
   - Subscribe tenant to a specific operation/API
   Body:
   - operationId
   Effects:
   - Creates an operation subscription with status ACTIVE

8. `DELETE /subscriptions/{connectorId}`
   Purpose:
   - Unsubscribe tenant from a connector (and its operations)
   Effects:
   - Marks subscription INACTIVE (soft delete)

9. `DELETE /subscriptions/{connectorId}/operations/{operationId}`
   Purpose:
   - Unsubscribe tenant from a specific operation
   Effects:
   - Marks operation subscription INACTIVE

### 3) Configuration + mapping

10. `GET /subscriptions/{connectorId}/operations/{operationId}/mapping`
    Purpose:
    - Get tenant’s selected fields and mapping rules for a specific operation

11. `PUT /subscriptions/{connectorId}/operations/{operationId}/mapping`
    Purpose:
    - Save tenant’s selected fields and mapping rules
    Body:
    - selectedFields
    - inputMapping (source -> target field mapping)
    - outputMapping
    - defaults

12. `GET /subscriptions/{connectorId}/operations/{operationId}/config`
    Purpose:
    - Get per-tenant operational config (timeouts, retries, rate limits)

13. `PUT /subscriptions/{connectorId}/operations/{operationId}/config`
    Purpose:
    - Save per-tenant operational config

### 4) Credentials

14. `PUT /subscriptions/{connectorId}/credentials`
    Purpose:
    - Store or update tenant credentials
    Body:
    - authType
    - credential payload (stored in Secrets Manager)

15. `POST /subscriptions/{connectorId}/test-connection`
    Purpose:
    - Validate credentials against connector (optional)

---

## Data model (DynamoDB)

The following tables support the control plane. All tables use soft-delete
(status) to preserve audit history.

### 1) Tenants
Table: `tenants`
- PK: `tenantId`
- Attributes: name, status, plan, createdAt, limits

Purpose:
- Stores tenant identity and plan/limits.

### 2) Connectors (global catalog)
Table: `connectors`
- PK: `connectorId`
- Attributes: name, version, status, description, supportedAuthTypes, tags

Purpose:
- Global connector catalog (shared by all tenants).

### 3) ConnectorOperations (global catalog)
Table: `connector_operations`
- PK: `connectorId`
- SK: `operationId`
- Attributes:
  - method, path, summary, tags
  - requestSchema (JSON schema)
  - responseSchema (JSON schema)
  - examples

Purpose:
- Global operation definitions tied to a connector.

### 4) Subscriptions (per tenant, per connector)
Table: `subscriptions`
- PK: `tenantId`
- SK: `connectorId`
- Attributes: status, createdAt, updatedAt, limitsOverride, configRef

Purpose:
- Stores which connectors a tenant has subscribed to.

### 5) OperationSubscriptions (per tenant, per operation)
Table: `operation_subscriptions`
- PK: `tenantId`
- SK: `connectorId#operationId`
- Attributes: status, createdAt, updatedAt, mappingRef, configRef

Purpose:
- Stores which operations a tenant has enabled.

### 6) OperationMappings (per tenant, per operation)
Table: `operation_mappings`
- PK: `tenantId`
- SK: `connectorId#operationId`
- Attributes:
  - selectedFields
  - inputMapping
  - outputMapping
  - defaults
  - validationRules

Purpose:
- Stores tenant-specific field selections and mapping rules.

### 7) ConnectorCredentials (per tenant, per connector)
Table: `connector_credentials`
- PK: `tenantId`
- SK: `connectorId`
- Attributes:
  - authType
  - secretRef (Secrets Manager ARN)
  - metadata (non-secret)

Purpose:
- Stores a reference to tenant credentials, never secrets directly.

---

## Best practices

Schema management
- Schemas are global and versioned. Tenants can only select fields, not modify schemas.
- Operation versions should be explicit (e.g., `dvla.lookup.v1`).

Secrets
- Store secrets in AWS Secrets Manager.
- DynamoDB stores only secret references.

Soft deletes
- Keep status fields (ACTIVE/INACTIVE) for auditability.

Audit + change tracking
- Store `createdAt`, `updatedAt`, `updatedBy` on subscription/mapping records.

---

## Next implementation steps

1) Add DynamoDB table definitions in Terraform.
2) Add REST endpoints in Hub Service (controllers + DTOs).
3) Add repository/service layer for each table.
4) Add Secrets Manager integration for credentials.
5) Provide example connector catalog seed data.






## Global API schema storage

Source of truth
- Store canonical schemas in the repo as versioned files.
- Suggested layout: `schemas/{connectorId}/{operationId}.json`.

Runtime storage
- Store a copy of the schema in DynamoDB (`connector_operations` table):
  - `requestSchema` and `responseSchema` JSON blobs.
- Keep `method` and `path` as separate attributes for fast UI listing.

Format
- Use JSON Schema (draft 2020-12) for request/response bodies.
- Optionally store OpenAPI fragments, but keep schemas as JSON Schema blocks.

Sync flow
1) Update schema file in repo.
2) Seed or sync into DynamoDB (seed script or CI job).
3) Hub APIs read from DynamoDB for UI rendering.

---

## Permission model (tenant roles)

Goal
- Allow admins to manage subscriptions and credentials.
- Allow read-only users to view connectors, specs, and observability.

Where permissions are defined
- Define roles and permissions in the Hub Service (control plane).
- Store per-user role assignments in a tenant-scoped table.

Suggested roles
- `TENANT_ADMIN`: full access (create/update/delete subscriptions, credentials, mappings).
- `TENANT_EDITOR`: modify mappings/config; no credential changes.
- `TENANT_VIEWER`: read-only access to catalog, subscriptions, and observability.

Data model additions
Table: `tenant_user_roles`
- PK: `tenantId`
- SK: `userId`
- Attributes: role, status, createdAt, updatedAt

Enforcement in Hub Service
- JWT contains user identity and tenant context.
- Hub Service maps `userId -> role` via `tenant_user_roles`.
- Each endpoint checks role requirements:
  - Read endpoints: allow all roles.
  - Write endpoints: admin/editor only.
  - Credentials endpoints: admin only.

UI control
- UI calls `/me/permissions` (or similar) to retrieve role/permissions.
- UI shows/hides actions based on permissions:
  - Admin: manage connectors, credentials, subscriptions.
  - Editor: manage mappings/config.
  - Viewer: read-only screens.
- Backend still enforces permissions; UI is only for UX.



## Authorization control (Cognito + DynamoDB)

Design principle
- Cognito handles identity and login.
- DynamoDB remains the source of truth for tenant‑scoped roles and permissions.

Why
- Cognito groups are global per user pool and do not model per‑tenant roles well.
- DynamoDB allows one user to have different roles across tenants.

Recommended model
Table: `tenant_user_roles`
- PK: `tenantId`
- SK: `userId` (Cognito `sub`)
- Attributes: role, status, createdAt, updatedAt

Authorization enforcement
- Hub Service resolves `tenantId` and `userId` from the JWT.
- Hub Service fetches role from DynamoDB (or caches it).
- API layer enforces role requirements per endpoint.

Optional optimization
- Use Cognito Pre Token Generation Lambda to add custom claims:
  - `tenant_id`, `role`
- Hub Service can trust claims for faster authorization, with periodic validation.

UI usage
- UI calls a Hub endpoint such as `GET /me/permissions` and
  renders controls accordingly.
- Backend enforcement is still mandatory.






