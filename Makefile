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
sqlcyaml:
	sqlc init
generatesqlcfromyaml:
	sqlc generate
installpgengine:
	go get github.com/lib/pg
test:
	go test -v -cover ./...
server:
	go run main.go
install_go-gin:
	go get -u github.com/gin-gonic/gin
installviper_env:
	go get github.com/spf13/viper
installgomock_mockdb4testing:
	go get github.com/golang/mock/mockgen@v1.6.0
mockdb:
	mockgen -package mockdb -destination db/mock/store.go github.com/muhsufyan/transaksi_transfer/db/sqlc Store
.PHONY: pull_postgres12alpine new_container_postgres installsqlc run_postgres createdb migratesqlc installgolangmigrate dropdb migrateup migratedown sqlcyaml generatesqlcfromyaml installpgengine test server install_go-gin installviper_env installgomock_mockdb4testing mockdb
