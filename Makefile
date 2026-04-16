APP=plexus-bff-service

.PHONY: test build run fmt

test:
	go test ./...

build:
	go build -o bin/$(APP) ./cmd/$(APP)

run:
	go run ./cmd/$(APP)

fmt:
	go fmt ./...
