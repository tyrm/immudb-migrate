PROJECT_NAME=immudb-migrate

fmt:
	@echo formatting
	@go fmt $(shell go list ./... | grep -v /vendor/)

test:  tidy fmt
	go test -race -cover ./...

tidy:
	go mod tidy -compat=1.17

.PHONY: fmt test tidy
