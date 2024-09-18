.PHONY=build
# include .env
# export $(shell sed 's/=.*//' .env)
build-client:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/grpc-client cmd/client/main.go

build-server:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/grpc-server cmd/server/main.go

run-client: build-client
	@./bin/grpc-client

run-server: build-server
	@./bin/grpc-server

generate-client:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
coverage:
	@go test -v -cover ./...

test:
	@go test -v ./...
