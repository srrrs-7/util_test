.PHONY: gopher rust node deno bun php linux k6 plantuml

gopher:
	docker compose build gopher
	docker compose run --rm gopher bash
rust:
	docker compose build rust
	docker compose run --rm rust bash
node:
	docker compose build node
	docker compose run --rm node bash
deno:
	docker compose build deno
	docker compose run --rm deno bash
bun:
	docker compose up -d bun --build
deno-server:
	docker compose up -d deno --build
php:
	docker compose build php
	docker compose run --rm php bash
linux:
	docker compose build linux
	docker compose run --rm linux bash
k6:
	docker compose run --rm k6 run --vus 10 --duration 5s /k6/script.js
plantuml:
	docker compose up -d plantuml