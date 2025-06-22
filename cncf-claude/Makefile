.PHONY: up build down gopher grpc claude-master claude-worker

up:
	docker compose up -d claude-master claude-worker queue --scale claude-worker=1

build:
	docker compose build
	docker image ls | grep "<none>" | awk '{print $ 3}' | xargs docker rmi

down:
	docker compose down

gopher:
	docker compose run --rm gopher sh

grpc:
	docker compose run --rm gopher protoc --proto_path=/go/src/driver/grpc/proto \
		--go_out=. --go-grpc_out=. /go/src/driver/grpc/proto/queue.proto

# exec claude-master
claude-master:
	docker compose run --rm claude-master

claude-worker:
	docker compose exec claude-worker sh