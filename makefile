.DEFAULT_GOAL := build

.PHONY:goimports gofumpt staticcheck test build clean

goimports:
	goimports -l -w .

gofumpt: goimports
	gofumpt -l -w .

staticcheck: gofumpt
	staticcheck ./...

test: staticcheck
	go test ./... -v -cover
	
build: test
	go build -C cmd/web -o ../../bin/blackdog-web

clean: build
	go clean

