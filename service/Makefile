.PHONY: protos server streams queues

protos:
	protoc --go_out=./server/namespaces --go_opt=paths=source_relative --go-grpc_out=./server/namespaces --go-grpc_opt=paths=source_relative protos/namespaces.proto
	protoc --go_out=./server/pages --go_opt=paths=source_relative --go-grpc_out=./server/pages --go-grpc_opt=paths=source_relative protos/pages.proto
	protoc --go_out=./server/projects --go_opt=paths=source_relative --go-grpc_out=./server/projects --go-grpc_opt=paths=source_relative protos/projects.proto
	protoc --go_out=./server/search --go_opt=paths=source_relative --go-grpc_out=./server/search --go-grpc_opt=paths=source_relative protos/search.proto

server:
	docker-compose up -d --no-deps --build server

streams:
	docker-compose up -d --no-deps --build streams

queues:
	docker-compose up -d --no-deps --build queues
