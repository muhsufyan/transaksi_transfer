package main

import (
	// tanpa ini _"github.com/lib/pq" we cant talk to db
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/muhsufyan/transaksi_transfer/api"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
)

/*
	to create server first need connect to db & create Store
	codenya sama sprti di db/sqlc/main_test.go bagian const
*/
const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:password@localhost:5432/bank?sslmode=disable"
	serverAddress = "0.0.0.0:8000"
)

func main() {
	// konek to db
	koneksi, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("tdk tersambung ke db karena error :", err)
	}
	store := db.NewStore(koneksi)
	server := api.NewServer(store)
	// to start the server
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cant start server :", err)
	}
}
