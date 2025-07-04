.PHONY: up build down gopher grpc tidy

up:
	docker compose up -d client server --scale client=1

build:
	docker compose build
	docker image ls | grep "<none>" | awk '{print $ 3}' | xargs docker rmi

down:
	docker compose down

gopher:
	docker compose run --rm gopher fish

grpc:
	docker compose run --rm gopher protoc --proto_path=/go/src/driver/grpc/proto \
		--go_out=. --go-grpc_out=. /go/src/driver/grpc/proto/chat.proto

tidy:
	docker compose run --rm gopher go mod tidy