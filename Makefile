TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go' | grep -vE './_local')
GO_CMD ?= go
BUILD_DIR = $(PWD)/dist
SHELL := /bin/bash

all: clean tidy fmt lint security test

setup:
	@command -v golangci-lint 2>&1 > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@command -v gosec 2>&1 > /dev/null || go install github.com/securego/gosec/v2/cmd/gosec@latest
	@command -v goreleaser 2>&1 > /dev/null || go install github.com/goreleaser/goreleaser@latest

clean:
	rm -rf ./dist

tidy:
	go mod tidy

fmt:
	$(GO_CMD)fmt -w $(GOFMT_FILES)

lint:
	golangci-lint run

security:
	gosec -exclude-dir _local -quiet ./...

test:
	go test -v -timeout 30s -coverprofile=cover.out -cover $(TEST)
	go tool cover -func=cover.out
