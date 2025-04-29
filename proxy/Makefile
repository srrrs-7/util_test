.PHONY: up down client mysql proxy logs

up:
	docker compose up -d --build
	docker images -q --filter "dangling=true" | xargs docker rmi --force

down:
	docker compose down --remove-orphans --volumes

client:
	docker compose run --rm client /usr/local/bin/client

mysql-client:
	docker compose run --rm client mysql -h proxy -P 8080 -u root -p

mysql:
	docker compose exec -it mysql mysql -u root -p

proxy:
	docker compose exec -it proxy sh

logs:
	docker compose logs -f --tail 100