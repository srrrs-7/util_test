.PHONY: init env test test-f migrate rollback migrate-f request model controller provider clear

init:
	cp .env.example .env

env:
	cd ../ && docker compose exec laravel-web php artisan env

test:
	cd ../ && docker compose exec laravel-web php artisan test --parallel --processes=4

test-f:
	cd ../ && docker compose exec laravel-web php artisan make:test

migrate:
	cd ../ && docker compose exec laravel-web php artisan migrate

rollback:
	cd ../ && docker compose exec laravel-web php artisan migrate:rollback

migrate-f:
	cd ../ && docker compose exec laravel-web php artisan make:migration

request:
	cd ../ && docker compose exec laravel-web php artisan make:request

model:
	cd ../ && docker compose exec laravel-web php artisan make:model

controller:
	cd ../ && docker compose exec laravel-web php artisan make:controller

provider:
	cd ../ && docker compose exec laravel-web php artisan make:provider

clear:
	cd ../ && docker compose exec laravel-web sh -c "php artisan config:clear && php artisan route:clear && php artisan view:clear && php artisan cache:clear"

# frontend
.PHONY: node-b
node-b:
	cd ../ && docker compose up -d laravel-web --build