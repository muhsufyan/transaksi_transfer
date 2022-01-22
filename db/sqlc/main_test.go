package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/bank?sslmode=disable"
)

// koneksi ke db lihat db.go struct Queries
var testQueries *Queries

// buat func dg param testing.T (setiap testing hrs ada param tipe param *testing.T)
func TestMain(m *testing.M) {
	// koneksi ke db
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("tdk tersambung ke db karena error :", err)
	}
	// jika tdk ada error, konek u/ membuat objek testQueries yg baru. New() dibuat di db.go (hsl generate)
	testQueries = New(conn)

	//running unit test
	os.Exit(m.Run())
}
