
SHELL=/bin/bash
GO_ROOT := $(shell go env GOROOT)
export PATH := $(GO_ROOT)/bin:$(PATH)


.PHONY: install
install:
	go install go.osspkg.com/goppy/v3/cmd/goppy@latest
	goppy setup-lib

.PHONY: lint
lint:
	goppy lint

.PHONY: generate
generate:
	go generate ./...

.PHONY: license
license:
	goppy license

.PHONY: build
build:
	goppy build --arch=amd64

.PHONY: tests
tests:
	go test ./...

.PHONY: verify
verify: generate lint tests
	go test -race ./...
	go vet ./...
	go mod verify
	git diff --check

.PHONY: pre-commit
pre-commit: install license generate lint tests build

.PHONY: ci
ci: verify build

.PHONY: vulncheck
vulncheck:
	govulncheck ./...
