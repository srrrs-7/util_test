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

.PHONY: rds mongo goapi goworker gopher rust rustapi rustworker node deno bun php linux k6 plantuml
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
rustapi:
	docker compose up -d rustapi --build
rustworker:
	docker compose up -d rustworker --build
node:
	docker compose build node
	docker compose run --rm node bash
deno:
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

.PHONY: k3s k3s-get-pods k3s-describe-pod k3s-mysql k3s-mysql-conn k3s-postgres k3s-postgres-conn
k3s:
	docker compose up -d k3s
k3s-get-pods:
	docker compose exec k3s kubectl get svc,sts,po 
k3s-describe-pod:
	docker compose exec k3s kubectl describe pod mysql
k3s-delete-pod:
	docker compose exec k3s kubectl delete pod ${POD}
k3s-delete-all:
	docker compose exec k3s kubectl get pods | awk '{print $1}' | grep -v NAME | xargs kubectl delete pods
k3s-mysql:
	docker compose exec k3s kubectl apply -f /k3s/mysql.yaml
k3s-mysql-conn:
	docker compose exec k3s kubectl exec -it mysql -c mysql -- mysql -u root -p"password"
k3s-postgres:
	docker compose exec k3s kubectl apply -f /k3s/postgres.yaml
k3s-postgres-conn:
	docker compose exec k3s kubectl exec -it postgres -c mysql -- mysql -u root -p"password"

.PHONY: k8s-get-pods k8s-del-pod k8s-nginx k8s-mysql k8s-del-mysql k8s-postgres k8s-del-postgres
k8s-get-pods:
	kubectl get svc,sts,po 
k8s-del-pod:
	kubectl delete pods ${POD}
k8s-nginx:
	kubectl run nginx --image nginx:latest
	kubectl expose pod nginx --type=NodePort --port=80
k8s-mysql:
	kubectl apply -f ./k8s/mysql.yaml
	kubectl get pvc
	kubectl get pv
	kubectl get pods | grep mysql | awk '{print $$1}' | xargs kubectl describe pod
k8s-del-mysql:
	kubectl delete deployments mysql
	kubectl delete service mysql
	kubectl delete pvc mysql-pvc
	kubectl delete pv mysql-pv
k8s-postgres:
	kubectl apply -f ./k8s/postgres.yaml
	kubectl get pvc
	kubectl get pv
	kubectl get pods | grep postgres | awk '{print $$1}' | xargs kubectl describe pod
k8s-del-postgres:
	kubectl delete deployments postgresql
	kubectl delete service postgresql
	kubectl delete pvc postgres-pvc
	kubectl delete pv postgres-pv

.PHONY: orb-ssh
orb-ssh:
	ssh orb

.PHONY: laravel
laravel:
	docker compose up -d laravel-web laravel-db laravel-nginx --build