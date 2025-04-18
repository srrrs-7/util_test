.PHONY: init fmt stan phpmd phpcs env test test-f migrate rollback migrate-f request model controller provider clear

init:
	cp .env.example .env

fmt:
	cd .. && docker compose run --rm laravel-web php -d memory_limit=1024M vendor/bin/php-cs-fixer fix ./app
 
stan:
	cd .. && docker compose run --rm laravel-web php -d memory_limit=1024M vendor/bin/phpstan analyse ./app
 
phpmd:
	cd .. && docker compose run --rm laravel-web php -d memory_limit=1024M vendor/bin/phpmd ./app text cleancode
 
phpcs:
	cd .. && docker compose run --rm laravel-web php -d memory_limit=1024M vendor/bin/phpcs --standard=PSR12 ./app

rector:
	cd .. && docker compose run --rm laravel-web php -d memory_limit=1024M vendor/bin/rector process --dry-run ./app
	@read -p "Do you want to apply the changes? [y/N]: " answer; \
	if [ "$$answer" != "y" ]; then \
		exit 0; \
	fi
	cd .. && docker compose run --rm laravel-web php -d memory_limit=1024M vendor/bin/rector process ./app

env:
	cd ../ && docker compose run --rm laravel-web php artisan env

test:
	cd ../ && docker compose run --rm laravel-web php artisan test --parallel --processes=4

test-f:
	cd ../ && docker compose run --rm laravel-web php artisan make:test

migrate:
	cd ../ && docker compose run --rm laravel-web php artisan migrate

rollback:
	cd ../ && docker compose run --rm laravel-web php artisan migrate:rollback

migrate-f:
	cd ../ && docker compose run --rm laravel-web php artisan make:migration

request:
	cd ../ && docker compose run --rm laravel-web php artisan make:request

model:
	cd ../ && docker compose run --rm laravel-web php artisan make:model

controller:
	cd ../ && docker compose run --rm laravel-web php artisan make:controller

provider:
	cd ../ && docker compose run --rm laravel-web php artisan make:provider

job:
	cd ../ && docker compose run --rm laravel-web php artisan make:job

clear:
	cd ../ && docker compose run --rm laravel-web sh -c "php artisan config:clear && php artisan route:clear && php artisan view:clear && php artisan cache:clear"

# frontend
.PHONY: node-b
node-b:
	cd ../ && docker compose up -d laravel-web --build