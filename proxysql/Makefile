.PHONY: up down master slave proxysql proxy clean

up:
	docker compose up -d --build master slave proxysql

down:
	docker compose down

master:
	docker compose exec master mysql -uroot -proot

slave:
	docker compose exec slave mysql -uroot -proot

proxysql:
	docker compose exec proxysql mysql -h 127.0.0.1 -P 6032 -u admin -padmin

proxy:
	docker compose exec proxysql mysql -h 127.0.0.1 -P 6033 -u root -proot

clean:
	docker compose down --volumes
	rm -rf master
	rm -rf slave
	rm -rf proxysql