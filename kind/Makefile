# cluster commands for kind using Docker Compose
.PHONY: cluster-up cluster-build cluster-down cluster-logs

cluster-up:
	docker compose up -d

cluster-build:
	docker compose build

cluster-down:
	docker compose down

cluster-logs:
	docker compose logs -f

# Commands to interact with the cluster
.PHONY: cluster delete exec get-all

cluster:
	docker compose exec cluster-1 kind create cluster --config /src/cluster/config.yaml

delete:
	docker compose exec cluster-1 kind delete cluster

exec:
	docker compose exec cluster-1 sh

get-all:
	docker compose exec cluster-1 kubectl get all -A

# about node commands
.PHONY: get-nodes get-nodes-detail

get-nodes:
	docker compose exec cluster-1 kubectl get nodes -o json

get-nodes-detail:
	docker compose exec cluster-1 kubectl get nodes -o wide

# about pods commands
.PHONY: get-pods get-pods-detail apply-pod describe-deploy describe-pod top-pods

get-pods:
	docker compose exec cluster-1 kubectl get pods -o json

get-pods-detail:
	docker compose exec cluster-1 kubectl get pods -o wide

apply-pod:
	docker compose exec cluster-1 kubectl apply -f /src/pods/nginx.yaml

describe-deploy:
	docker compose exec cluster-1 kubectl describe deployment nginx-deployment

describe-pod:
	docker compose exec cluster-1 kubectl describe pod -l app=nginx

top-pods:
	docker compose exec cluster-1 kubectl top pods -l app=nginx
