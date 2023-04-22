SHELL := /bin/bash

.PHONY: build deps format run test

build: clean
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api cmd/api/api.go
	zip bin/api.zip bin/api
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/async cmd/sqs/sqs.go
	zip bin/sqs.zip bin/sqs
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/one cmd/one/one.go
	zip bin/one.zip bin/one
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/two cmd/two/two.go
	zip bin/two.zip bin/two
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/three cmd/three/three.go
	zip bin/three.zip bin/three

clean:
	rm -rfv ./bin

deploy-all: deploy-staging deploy-production

deploy-staging: pre-deploy
	serverless deploy --verbose --stage staging --region us-east-2 --org learnfullysean

deploy-production: pre-deploy
	serverless deploy --verbose  --stage production --region us-east-2 --org learnfullysean

deps:
	go install github.com/seantcanavan/fresh/v2@latest
	npm install -g serverless
	serverless plugin install -n serverless-lift

format:
	go fmt ./...

pre-deploy: clean build

run:
	source .env && fresh

test:
	go test ./...
