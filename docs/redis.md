# Zynchub Redis Usage

This document explains how Redis is used by the Hub Service, and how to configure it
for local and AWS environments.

Purpose
- Cache feature flags in user session scope (when enabled).
- Optional session storage and shared caching in dev/test/prod.

Profiles and behavior
- Local:
  - Default local profile uses in-memory/session stub and no auth.
  - Optional: `local-redis` profile uses embedded Redis to emulate server behavior.
- Dev/Test/Prod:
  - Uses AWS ElastiCache (Redis) with IAM authentication.

Configuration flags (application.yml)
- `features.session.enabled`
  - true: cache features in session (Redis-backed in dev/test/prod).
  - false: bypass session cache and always fetch from AppConfig.

Local setup
1) Run with `local-redis` profile:
   - Set `SPRING_PROFILES_ACTIVE=local-redis`
2) Embedded Redis starts automatically; no AWS connectivity required.

AWS setup (ElastiCache)
Terraform creates:
- ElastiCache replication group (Redis)
- ElastiCache user and user group (IAM auth)
- Security group to allow ECS -> Redis

ECS environment variables
These are set in the task definition (dev/test/prod):
- `SPRING_DATA_REDIS_HOST` (endpoint)
- `SPRING_DATA_REDIS_PORT` (6379)
- `SPRING_DATA_REDIS_SSL` (true)
- `SPRING_DATA_REDIS_USERID` (ElastiCache user ID)
- `SPRING_DATA_REDIS_REPLICATIONGROUPID` (replication group ID)

IAM permissions
ECS task role must allow:
- `elasticache:Connect` on the ElastiCache user and replication group

Parameter Store / Secrets
No Redis credentials are stored. IAM auth tokens are generated at runtime, so no
password or secret is needed.

Verification
- Confirm the ECS task definition contains all Redis env vars.
- Confirm the ECS task role has `elasticache:Connect` to the created user/group.
- Check CloudWatch logs for successful Redis connection at startup.

Troubleshooting
- WRONGPASS:
  - Usually indicates the task is using the wrong Redis user or old task definition.
  - Force a new ECS deployment after terraform apply.
- DNS or connection errors:
  - Validate VPC/subnet routing and security group rules.
