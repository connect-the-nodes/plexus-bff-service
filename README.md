# plexus-hub-service

Primary backend hub microservice that exposes REST (and optional GraphQL) endpoints for the frontend applications.

## Overview

The **plexus-hub-service** acts as the central control-plane service for the ZyncHub integration platform. It provides secure, versioned APIs consumed by React-based frontend applications and automation tools to configure, manage, and govern integrations with third-party data providers.

This service enables users to define integration contracts, canonical models, provider mappings, and integration workflows without requiring custom integration code. It is designed for enterprise insurance environments, supporting UK and EU regulatory requirements.

---

## Responsibilities

- Expose REST and GraphQL APIs for frontend applications
- Manage integration definitions and lifecycle
- Maintain versioned input and output contracts
- Manage canonical data models
- Configure provider metadata and routing rules
- Define and store field mappings
- Initiate asynchronous SDK generation workflows
- Provide auditability and traceability for configuration changes

---

## Key Concepts

### Integrations
Logical representations of a business capability (e.g. vehicle lookup, address lookup) that connect client systems to one or more external providers. Integrations are versioned and lifecycle-managed.

### Contracts
Versioned and immutable definitions of the data exchanged between client systems and the platform, including validation rules and error models.

### Canonical Models
Provider-agnostic representations of business entities used internally to normalize data across multiple third-party providers.

### Providers
External systems or services (real or stubbed) that supply data. Provider credentials are managed securely by the platform.

### Mappings
Declarative configurations that map fields between client contracts, canonical models, and provider-specific schemas.

### SDK Generation
Asynchronous workflows that generate client SDKs based on published contracts, including automated validation and testing.

---

## API Characteristics

- Contract-first design using OpenAPI
- OAuth2-secured endpoints
- Versioned APIs
- Designed for frontend consumption
- Supports synchronous and asynchronous workflows

---

## Security & Compliance

- OAuth2 authentication and authorization
- Fine-grained access control
- Correlation ID propagation
- Audit logging of all configuration changes
- PII-aware field handling
- Designed for UK and EU regulatory environments

---

## Architecture Context

The hub service operates as part of a broader integration platform:

- Frontend applications use this service to configure integrations
- Runtime execution is handled by downstream orchestration and provider connector services
- The service is stateless and cloud-native, designed for horizontal scalability

---

## Non-Goals

- Direct invocation of third-party provider APIs
- Runtime data processing or enrichment
- Long-running orchestration execution

---

## Documentation

- OpenAPI specifications are available in the `/openapi` directory
- Authentication and authorization details are defined in the security documentation
- Additional docs:
  - `docs/oauth-flow.md` (Cognito OAuth/OIDC flow)
  - `docs/redis.md` (Redis usage and configuration)
  - `docs/feature-flags.md` (Feature flagging approach)

---

## Local Redis (embedded)

For local development with session data persisted in Redis, use the embedded Redis profile:

- Profile: `local-redis`
- Redis host/port: `localhost:6380`
- Session store: `spring.session.store-type=redis`
- Namespace: `local:spring:session`

Example:

```bash
SPRING_PROFILES_ACTIVE=local-redis
```

The default `local` profile keeps `spring.session.store-type=none` and disables Redis auto-configuration.

---

## Dev Redis (AWS)

In `dev`, Redis sessions are backed by AWS ElastiCache (Redis) with IAM authentication enabled.
Terraform injects the required environment variables into the ECS task definition.

Required env vars (set by Terraform):

- `SPRING_DATA_REDIS_HOST` (ElastiCache primary endpoint)
- `SPRING_DATA_REDIS_PORT` (default `6379`)
- `SPRING_DATA_REDIS_SSL_ENABLED=true`
- `SPRING_DATA_REDIS_IAM_ENABLED=true`
- `SPRING_DATA_REDIS_USERID` (ElastiCache IAM user)
- `SPRING_DATA_REDIS_REPLICATIONGROUPID`
- `SPRING_DATA_REDIS_REGION` (e.g., `eu-west-1`)

The `dev` profile enables:
- `spring.session.store-type=redis`
- `spring.redis.namespace=dev:spring:session`

Verification:

- Redis health (when enabled): `GET /actuator/health` should include `"redis": {"status": "UP"}`.
- Session persistence: create a session via `POST /_test/session` (enabled in `dev` and `local-redis`) and confirm Redis keys are created under the namespace prefix.

Example:

```bash
curl -X POST http://localhost:8080/_test/session
```

---

## GitHub Actions deployment permissions

If the deploy workflow fails at the ECR login step with `ecr:GetAuthorizationToken`, the
OIDC role used by GitHub Actions is missing required ECR permissions. Attach a policy
similar to the following to the role referenced by `AWS_ROLE_ARN` (adjust `Resource`
scopes as needed for your account):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "EcrLogin",
      "Effect": "Allow",
      "Action": ["ecr:GetAuthorizationToken"],
      "Resource": "*"
    },
    {
      "Sid": "EcrPushPullRepo",
      "Effect": "Allow",
      "Action": [
        "ecr:BatchCheckLayerAvailability",
        "ecr:BatchGetImage",
        "ecr:CompleteLayerUpload",
        "ecr:DescribeImages",
        "ecr:DescribeRepositories",
        "ecr:GetDownloadUrlForLayer",
        "ecr:InitiateLayerUpload",
        "ecr:PutImage",
        "ecr:UploadLayerPart"
      ],
      "Resource": "arn:aws:ecr:<region>:<account-id>:repository/plexus-hub-*"
    }
  ]
}
```

Infrastructure code for Hub is now maintained centrally in:

`C:\Naresh\zynch_app\zynchub-platform-infra`

Relevant stack path:

`live/dev/eu-west-1/plexus/hub`

Observability integration configuration is managed there via:

- Terraform variable: `observability_service_base_url`
- ECS env var: `OBSERVABILITY_SERVICE_BASE_URL`
