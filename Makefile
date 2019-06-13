SHELL=/bin/bash

PROJECT=c2ae

GIT_COMMIT=$(shell git rev-list -1 HEAD)
GIT_TAG=$(shell git describe --exact-match HEAD 2>/dev/null || true)
NOW=$(shell  date "+%Y%m%d")

GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH=amd64

C2AETEST_POSTGRES="${C2AETEST_POSTGRES:-}"

.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: test build-cli build-api ## Build the binaries

build-cli:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -race -o ./bin/${PROJECT}-cli -ldflags "-X main.gitTag=${GIT_TAG} -X main.gitCommit=${GIT_COMMIT} -X main.buildDate=${NOW}" ./cmd/cli

build-api:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -race -o ./bin/${PROJECT}-api -ldflags "-X main.gitTag=${GIT_TAG} -X main.gitCommit=${GIT_COMMIT} -X main.buildDate=${NOW}" ./cmd/api

.PHONY: test
test: ## Run tests
	@if ! test -z "$$C2AETEST_POSTGRES"; then echo "C2AETEST_POSTGRES => enabled"; else echo "C2AETEST_POSTGRES => disabled"; fi; \
	go test -v -coverprofile=/tmp/go-code-cover -race -timeout 10s  ./...

.PHONY: generate
generate: ## Generate mocks and proto files
	go generate ./...
	protoc --proto_path . api.proto --go_out=plugins=grpc:./internal/pb

.PHONY: cover
cover: test ## Show coverage
	go tool cover -html=/tmp/go-code-cover

.PHONY: clean
clean: ## Restore project to pristine state
	rm -f bin/*
