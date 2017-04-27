
PACKAGE = github.com/mikattack/mlog
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`

.PHONY: vendor check fmt lint test test-cover-html help
.DEFAULT_GOAL := help

check: test fmt vet

vendor:
	go get github.com/stretchr/testify

test: vendor
	go test

fmt:
	gofmt -l *.go

vet:
	go vet

coverage:
	# echo "mode: count" > coverage-all.out
	go test -coverprofile=coverage.out -covermode=count $(pkg);
	# tail -n +2 coverage.out >> coverage-all.out;
	go tool cover -html=coverage.out

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
