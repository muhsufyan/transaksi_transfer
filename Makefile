pull_postgres12alpine:
	docker pull postgres:12-alpine
new_container_postgres:
	docker run --name postgres12alpine -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine
installsqlc:
	snap install sqlc
run_postgres:
	docker start postgres12alpine
migratesqlc:
	migrate create -ext sql -dir db/migration -seq init_schema
installgolangmigrate:
	curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
	apt-get update
	apt-get install -y migrate
createdb:
	docker exec -it postgres12alpine createdb --username=root --owner=root bank
dropdb:
	docker exec -it postgres12alpine dropdb bank 
migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose down
.PHONY: pull_postgres12alpine new_container_postgres installsqlc run_postgres createdb migratesqlc installgolangmigrate dropdb migrateup migratedown

https://www.youtube.com/watch?v=6_CH8Gx414A docker-compose