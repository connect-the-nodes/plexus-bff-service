# Parity Gap Report

## Verified in this workspace

- Repository structure created as a sibling of the Java service
- Java service behavior, configuration, and tests were mapped into the Go implementation plan
- Go dependencies resolved with `go mod tidy`
- Go test suite passed with `go test ./...`
- Go binary built successfully with `go build ./cmd/plexus-bff-service`
- Full repository compile passed with `go build ./...`
- Local runtime smoke test passed with `/actuator/health` returning HTTP `200`
- JWT/JWK validation support was implemented in the Go security layer and compiles successfully

## Not yet verified in this workspace

- runtime parity against real Redis
- runtime parity against AWS AppConfig
- runtime parity against real JWT/JWK verification
- container build execution
- execution of `internal/app/security` package tests is blocked by local Windows Application Control policy

## Known implementation gaps

- Redis-backed session persistence is implemented as a Go-native session store rather than Spring Session semantics; this preserves the development use case but still needs live verification for production parity.
- AWS Redis IAM token generation is documented for follow-up implementation if strict ElastiCache IAM parity is required beyond the current migration scaffold.
- JWT parsing and authority extraction are verified for the current local/test bearer-token flow, but real issuer/JWK execution tests on this host are currently blocked by Application Control policy.
