.PHONY: compile-pb

compile-pb:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/svcgrpc/cacheservice.proto