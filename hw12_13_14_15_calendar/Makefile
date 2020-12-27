BIN := "./bin/calendar"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

.PHONY: run
run: build
	$(BIN) -config ./configs/config.toml

.PHONY: version
version: build
	$(BIN) version

.PHONY: test
test:
	go test -v -race ./...

.PHONY: test-pq
test-pq:
	PQ_TEST="postgresql://calendar:calendar@localhost:5432/calendar" go test -v -count=1 -race -timeout=30s ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

.PHONY: lint
lint: install-lint-deps
	golangci-lint run ./...
