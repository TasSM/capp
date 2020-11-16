.PHONY: compile-pb unit-test test-client run build

compile-pb:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/svcgrpc/cacheservice.proto

unit-test:
	go test -race ./...

test-client:
	go run -race cmd/capp/test-client/main.go

run:
	go run -race cmd/capp/main.go

build:
	go build -o /bin/main cmd/capp/main.go