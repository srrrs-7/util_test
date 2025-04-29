.PHONY: up down vt-init

up:
	docker compose up -d --build db etcd vtctld vttablet1 vtgate 

down:
	docker compose down

vt-init:
	docker compose exec vtctld /vt/bin/vtctlclient -server vtctld:15999 CreateKeyspace test_keyspace
	docker compose exec vtctld /vt/bin/vtctlclient -server vtctld:15999 InitShardPrimary -force test_keyspace/0 test-0000000100

mysql:
	docker compose run --rm mysql-client bash -c "apt-get update \
		&& apt-get install -y default-mysql-client \
		&& mysql -h vigate -P 15001 -u root -p"