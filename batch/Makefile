.PHONY: protos listener server

protos:
	protoc --go_out=./server/diffs --go_opt=paths=source_relative --go-grpc_out=./server/diffs --go-grpc_opt=paths=source_relative protos/diffs.proto

listener:
	docker-compose up -d --no-deps --build listener

server:
	docker-compose up -d --no-deps --build server
