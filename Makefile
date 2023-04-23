SHELL := /bin/bash

.PHONY: build deps format run test

build: clean
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api cmd/api/api.go
	zip bin/api.zip bin/api
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/sqs cmd/sqs/sqs.go
	zip bin/sqs.zip bin/sqs
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/one cmd/1/one.go
	zip bin/one.zip bin/one
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/two cmd/2/two.go
	zip bin/two.zip bin/two
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/three cmd/3/three.go
	zip bin/three.zip bin/three

clean:
	rm -rfv ./bin

deploy-all: deploy-staging deploy-production

deploy-staging: pre-deploy
	serverless deploy --verbose --stage staging --region us-east-2 --org f72e1c13062e4f45ad951530acf9e5a7

deploy-production: pre-deploy
	serverless deploy --verbose  --stage production --region us-east-2 --org f72e1c13062e4f45ad951530acf9e5a7

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
