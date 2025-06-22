.PHONY: up build down gopher

up:
	docker compose up -d

build:
	docker compose build

down:
	docker compose down

gopher:
	docker compose run --rm gopher sh

grpc:
	docker compose run --rm gopher protoc --proto_path=/go/src/driver/grpc/proto \
		--go_out=. --go-grpc_out=. /go/src/driver/grpc/proto/queue.proto