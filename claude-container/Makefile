PHONY: up build ai

up:
	docker compose up -d ai

build:
	docker compose build
	docker image ls | grep "<none>" | awk '{print $ 3}' | xargs docker rmi

ai: up
	docker compose exec ai tmux
