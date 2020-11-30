PKGS := $(shell go list ./... | grep -v -E mocks)
SRCS := $(shell find . go.mod -name '*.go') go.mod
export GO111MODULE=on

VERSION := $(shell cat Version)

# Set the code coverage target for this project
COVER_TARGET ?= 85

LINT_OPTS ?= --fix
LINTERS=-E gofmt -E govet -E errcheck -E staticcheck -E gosimple -E structcheck \
        -E varcheck -E ineffassign -E typecheck -E unused

.PHONEY: build
build: ${SRCS} 
	@go build
.PHONY: tests
tests:  ## Run test suite
	@go test -race ${PKGS}

.PHONY: cover
cover:  ## Generate test coverage results
	@go test --covermode=count -coverprofile cover.profile ${PKGS} 
	@go tool cover -html cover.profile -o cover.html
	@go tool cover -func cover.profile -o cover.func
	@tail -n 1 cover.func | awk '{if (int($$3) > ${COVER_TARGET}) {print "Coverage good: " $$3} else {print "Coverage is less than ${COVER_TARGET}%: " $$3; exit 1}}'


.PHONY: lint
lint:  ## Run golint and go fmt on source base
	@golangci-lint run ${LINT_OPTS} --no-config --disable-all ${LINTERS} ./...

.PHONY: clean
clean:  ## Clean up any generated files
	@rm -f cover.*

.DEFAULT_GOAL := help
.PHONY: help
help:   ## Display this help message
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:  ## Show the version the Makefile will build
	@echo ${VERSION}

