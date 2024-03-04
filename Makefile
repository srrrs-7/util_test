.PHONY: gopher rust node deno php linux k6 plantuml

gopher:
	docker compose run --rm gopher bash
rust:
	docker compose run --rm rust bash
node:
	docker compose run --rm node bash
deno:
	docker compose run --rm deno bash
deno-server:
	docker compose up -d deno --build
php:
	docker compose run --rm php bash
linux:
	docker compose run --rm linux bash
k6:
	docker compose run --rm k6 run --vus 10 --duration 5s /k6/script.js
plantuml:
	docker compose up -d plantuml