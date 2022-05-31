PROJECT_NAME=immudb-migrate

check:
	golangci-lint run

check-fix:
	golangci-lint run --fix

fmt:
	@echo formatting
	@go fmt $(shell go list ./... | grep -v /vendor/)

test:  tidy fmt
	go test -race -cover ./...

tidy:
	go mod tidy -compat=1.17

.PHONY: check check-fix fmt test tidy
