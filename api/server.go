package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
)

// Server serves  HTTP request for banking service
type Server struct {
	// objek *sql.DB ada di db/sqlc/store.go bagian struct Store.
	// objek ini bertanggung jwb agar dpt terhub ke database ketika client melakukan request ke API
	store *db.Store
	// gin.Engin. send each request ke handler yg sesuai
	router *gin.Engine
}

// new HTTP server, all HTTP API route for service is here
func NewServer(store *db.Store) *Server {
	// new server
	server := &Server{store: store}
	// buat new router
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	// :id is url param
	router.GET("/account/:id", server.getAccount)
	// get list accounts with pagination
	router.GET("/account", server.listAccount)

	// route API new account
	// disini kita bisa masukkan banyak func sprti middleware, handler, dll. tp sekarang hanya handler saja
	// method ini adlh struct Server yg perlu we implement krn we mengakses objek store u/ menyimpan account baru ke db. implementnya ada di api/account.go
	// save objek router ke server.router
	server.router = router
	return server
}

// to run HTTP server on specific address
func (server *Server) Start(address string) error {
	// field route is private itulah alasannya kita buat Start is public
	return server.router.Run(address)
}

/*
gin.H is map[string]interface{} so we can store data key value apapun yg we want
*/
func errorResponse(err error) gin.H {
	// for now return error message
	return gin.H{"error": err.Error()}
}
