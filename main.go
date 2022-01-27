package main

import (
	// tanpa ini _"github.com/lib/pq" we cant talk to db
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/muhsufyan/transaksi_transfer/api"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
	"github.com/muhsufyan/transaksi_transfer/util"
)

/*
	to create server first need connect to db & create Store
	codenya sama sprti di db/sqlc/main_test.go bagian const
*/

func main() {
	// config lewat env
	config, err := util.LoadConfig(".") //"." karena main.go dan app.env ada di dir yg sama (root)
	if err != nil {
		log.Fatal("tdk bisa load config :", err)
	}
	// konek to db
	koneksi, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("tdk tersambung ke db karena error :", err)
	}
	store := db.NewStore(koneksi)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cant create server :", err)
	}
	// to start the server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cant start server :", err)
	}
}
