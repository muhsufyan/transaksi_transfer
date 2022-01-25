package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"github.com/muhsufyan/transaksi_transfer/util"
	_ "github.com/lib/pq"
)

// lihat db.go struct Queries
var testQueries *Queries

// koneksi ke db
var testDB *sql.DB

// buat func dg param testing.T (setiap testing hrs ada param tipe param *testing.T)
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..") // app.env ada di root dir sedangkan main_test.go ada di root/db/sqlc jd kita perlu ke root dir dg perintah "../.."
	if err != nil{
		log.Fatal("tdk bisa load config :", err)
	}
	// koneksi ke db
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("tdk tersambung ke db karena error :", err)
	}
	// jika tdk ada error, konek u/ membuat objek testQueries yg baru. New() dibuat di db.go (hsl generate)
	testQueries = New(testDB)

	//running unit test
	os.Exit(m.Run())
}
