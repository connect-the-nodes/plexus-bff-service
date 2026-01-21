# Zynchub Feature Flagging

This document explains how feature flags are managed and accessed.

Source of truth
- AWS AppConfig hosts the feature configuration (feature.yml).
- The Hub Service retrieves flags from AppConfig at runtime.

Retrieval flow
1) Hub Service requests the latest AppConfig configuration.
2) Response is parsed into feature objects.
3) If `features.session.enabled=true`, results are cached in session (Redis-backed in dev/test/prod).
4) Controller endpoints return the flags to the UI or other clients.

Configuration files
- `application.yml`:
  - Base defaults, including feature session toggle.
- `application-dev.yml` / `application-test.yml` / `application-prod.yml`:
  - `features.session.enabled=true` (default in environments).
- `application-local-redis.yml`:
  - Enables Redis session cache locally.

Key properties
- `features.session.enabled` (true/false)
  - true: cache feature flags in session.
  - false: always fetch directly from AppConfig.

AppConfig identifiers (AWS)
- Application name
- Environment name
- Configuration profile name

These are passed to ECS via environment variables:
- `APPCONFIG_APPLICATION`
- `APPCONFIG_ENVIRONMENT`
- `APPCONFIG_PROFILE`

Operational usage
- Update feature flags centrally in AppConfig (feature.yml).
- No redeploy required to pick up changes; AppConfig retrieval fetches current config.

Testing
- Local (no auth):
  - `GET /api/v1/features` returns current config (from local mock/feature file).
- Dev/Test/Prod:
  - Obtain Cognito access token, call API with `Authorization: Bearer <token>`.
  - Confirm flags are cached in Redis when `features.session.enabled=true`.

Disabling feature caching
- Set `features.session.enabled=false` in the target profile.
- Redeploy or restart to apply profile changes.
