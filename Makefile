.PHONY: surreal surql
surreal:
	docker compose up -d surreal surql --build
surql:
	docker compose run --rm surql

.PHONY: mysql data conn dump
mysql:
	docker compose up -d mysql --build
data:
	cd database/util && ./create_data.sh name age 1996-08-25 1000
conn:
	docker compose exec mysql sh -c "mysql -u root -p test < /mysql/init/init.sql"
	docker compose exec mysql mysql -u root -p test
dump:
	docker compose exec mysql mysqldump -u root -p test > ./database/dump/dump.sql
	gzip ./database/dump/dump.sql

.PHONY: postgres psql redash
postgres:
	docker compose up -d postgres --build
	docker compose exec postgres psql -h localhost -p 5432 -U root -d test -f /postgres/init/init.sql
psql:
	docker compose exec postgres psql -h localhost -p 5432 -U root -d test
redash: rds
	docker compose up -d redash redash-pg redash-scheduler redash-worker --build
	docker compose run --rm redash create_db

.PHONY: elasticsearch kibana
elasticsearch: kibana
	docker compose up -d elasticsearch --build
	cd ./elasticsearch && chmod +x ./init.sh && ./init.sh
kibana:
	docker compose up -d kibana --build

.PHONY: sqs send receive 
SQS_ENDPOINT="http://sqs:9324"
QUEUE_URL="http://sqs:9324/000000000000/test"
sqs:
	docker compose up -d sqs --build
	sleep 1
	docker compose run --rm aws sqs create-queue --queue-name=test --endpoint-url=$(SQS_ENDPOINT)
enqueue:
	docker compose run --rm aws sqs send-message --endpoint-url=$(SQS_ENDPOINT) --queue-url=${QUEUE_URL} --message-body="hello sqs"
dequeue:
	docker compose run --rm aws sqs receive-message --endpoint-url=$(SQS_ENDPOINT) --queue-url=${QUEUE_URL}

.PHONY: rds mongo goapi goworker gopher rust node deno bun php linux k6 plantuml
rds:
	docker compose up -d rds --build
mongo:
	docker compose up -d mongo --build
	docker compose exec mongo bash /mongo/init/mongo.sh
	docker compose exec mongo bash
goapi: mysql sqs rds goworker
	docker compose up -d goapi --build
goworker: mysql sqs rds goapi
	docker compose up -d goworker --build
worker-test:
	cd ./goapi/testdata && chmod +x ./worker-test.sh && ./worker-test.sh
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
deno-vite:
	docker compose run --rm deno deno run -A npm:create-vite-extra@latest
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
	docker compose run --rm k6 run /k6/scenario/scenario1/script.js --out json=/k6/log/out.json
plantuml:
	docker compose up -d plantuml

.PHONY: observer vector
observer:
	docker compose up -d prometheus grafana tempo loki otel-collector jaeger vector --build
vector:
	docker compose exec vector vector top

.PHONY: auth-doc master-doc stamp-doc shift-doc holiday-doc attendance-doc audit-doc
auth-doc:
	docker compose run --rm redoc npx @redocly/cli build-docs /app/auth.yaml --output /app/doc/auth.html
master-doc:
	docker compose run --rm redoc npx @redocly/cli build-docs /app/master.yaml --output /app/doc/master.html
stamp-doc:
	docker compose run --rm redoc npx @redocly/cli build-docs /app/stamp.yaml --output /app/doc/stamp.html
shift-doc:
	docker compose run --rm redoc npx @redocly/cli build-docs /app/shift.yaml --output /app/doc/shift.html
holiday-doc:
	docker compose run --rm redoc npx @redocly/cli build-docs /app/holiday.yaml --output /app/doc/holiday.html
attendance-doc:
	docker compose run --rm redoc npx @redocly/cli build-docs /app/aggregation.yaml --output /app/doc/aggregation.html
audit-doc:
	docker compose run --rm redoc npx @redocly/cli build-docs /app/audit.yaml --output /app/doc/audit.html