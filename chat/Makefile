.PHONY: client server build down gopher grpc tidy rmi logs

client:
	docker compose run --rm --service-ports client

server:
	docker compose up server

build:
	docker compose build

down:
	docker compose down

gopher:
	docker compose run --rm gopher fish

grpc:
	docker compose run --rm gopher protoc --proto_path=/go/src/driver/grpc/proto \
		--go_out=. --go-grpc_out=. /go/src/driver/grpc/proto/chat.proto

tidy:
	docker compose run --rm gopher go mod tidy

rmi:
	docker image ls | grep "<none>" | awk '{print $$3}' | xargs docker rmi

logs:
	docker compose logs