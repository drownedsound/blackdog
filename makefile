.DEFAULT_GOAL := build

# TODO: Create switch to run unit tests
# TODO: Create switch to run unit tests and compile code
.PHONY:goimports gofumpt staticcheck test build clean

goimports:
	goimports -l -w .

gofumpt: goimports
	gofumpt -l -w .

staticcheck: gofumpt
	staticcheck ./...

test: staticcheck
	go test ./internal/features/application ./internal/infra -cover

build: test
	go build -C cmd/web -o ../../bin/blackdog-web

clean: build
	go clean

