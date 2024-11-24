BINDIR ?= ./bin
CLI_BINARY = flightdeck

.PHONY: all
all: cli

.PHONY: cli
cli:
	GOARCH=amd64 GOOS=linux GOPRIVATE=github.com/arctir go build -ldflags="-X main.commit=$(shell git rev-parse HEAD)" -o ${BINDIR}/${CLI_BINARY}-linux-amd64 .
	GOARCH=amd64 GOOS=darwin GOPRIVATE=github.com/arctir go build -ldflags="-X main.commit=$(shell git rev-parse HEAD)" -o ${BINDIR}/${CLI_BINARY}-darwin-amd64 .
	GOARCH=amd64 GOOS=windows GOPRIVATE=github.com/arctir go build -ldflags="-X main.commit=$(shell git rev-parse HEAD)" -o ${BINDIR}/${CLI_BINARY}-windows-amd64 .
