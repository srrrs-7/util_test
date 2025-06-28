.PHONY: up build down ai

up:
	docker compose up -d ai

build:
	docker compose build ai
	sleep 1
	docker image ls | grep none | awk '{print $$3}' | xargs docker rmi

down:
	docker compose down

ai: up
	docker compose exec ai tmux