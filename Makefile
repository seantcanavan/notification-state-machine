SHELL := /bin/bash

.PHONY: build deps format run test

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/copilot main.go

clean:
	rm -rfv ./bin/

deps:
	go install github.com/seantcanavan/fresh/v2@latest

format:
	go fmt ./...

run:
	source .env && fresh

test:
	go test ./...
