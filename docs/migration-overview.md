# Migration Overview

This repository is a Go reimplementation of the Java `plexus-bff-service`.

## Migration principles

- Preserve externally observable behavior before refactoring internals.
- Keep package names close to the Java service layout for reviewability.
- Prefer explicit configuration over framework magic.
- Document every intentional deviation from the Java behavior.

## Current parity target

- REST endpoints and response shapes
- profile overlays under `configs/`
- feature loading from YAML and AWS AppConfig
- observability proxy endpoints
- auth redirect endpoints
- correlation ID echo/generation
- optional Redis connectivity and session write path

See `docs/parity-gap-report.md` for items that still require runtime verification.
