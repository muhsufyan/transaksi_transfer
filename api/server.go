package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
	"github.com/muhsufyan/transaksi_transfer/token"
	"github.com/muhsufyan/transaksi_transfer/util"
)

// Server serves  HTTP request for banking service
type Server struct {
	// objek *sql.DB ada di db/sqlc/store.go bagian struct Store.
	// objek ini bertanggung jwb agar dpt terhub ke database ketika client melakukan request ke API
	store db.Store //now to interface
	// gin.Engin. send each request ke handler yg sesuai
	router *gin.Engine
	// store token
	tokenMaker token.Maker
	// untuk ambil data config dr app.env
	config util.Config
}

// new HTTP server, all HTTP API route for service is here
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// buat obj token maker baru. kita bisa pilih jwt melalui token.NewJWTMaker() atau paseto melalui token.NewPasetoMaker(). Kali ini kita pilih paseto saja
	// perlu symmetris key as input/param
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey) //jika ingin mengganti jd jwt maka ubah kode ini jd token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	// new server
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
	// REGISTER CUSTOM VALIDATOR PD GIN. jika ok jlnkan RegisterValidation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// register custom validate func param 1 is tag, param 2 func validCurrency(our custom validate)
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

// routing method
func (server *Server) setupRouter() {
	// buat new router
	router := gin.Default()
	// buat route group yg menerapkan/menggunakan authMiddleware yg tlh kita buat
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// sekarang url yg menerapkan authMiddleware hrs melewati otorisasi dl
	authRoutes.POST("/accounts", server.createAccount)
	// :id is url param
	authRoutes.GET("/account/:id", server.getAccount)
	// get list accounts with pagination
	authRoutes.GET("/account", server.listAccount)
	// melakukan transfer baru
	authRoutes.POST("/transfers", server.createTransfer)
	// buat user baru
	router.POST("/users", server.createUser)
	// login request
	router.POST("/users/login", server.loginUser)
	// route API new account
	// disini kita bisa masukkan banyak func sprti middleware, handler, dll. tp sekarang hanya handler saja
	// method ini adlh struct Server yg perlu we implement krn we mengakses objek store u/ menyimpan account baru ke db. implementnya ada di api/account.go
	// save objek router ke server.router
	server.router = router
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
