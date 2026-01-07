# zynchub-digital-hub-service

Primary backend hub microservice that exposes REST (and optional GraphQL) endpoints for the frontend applications.

## Overview

The **zynchub-digital-hub-service** acts as the central control-plane service for the ZyncHub integration platform. It provides secure, versioned APIs consumed by React-based frontend applications and automation tools to configure, manage, and govern integrations with third-party data providers.

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

---
