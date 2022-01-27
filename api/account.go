package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
)

// for store the create account request
type createAccountRequest struct {
	// isinya sama dg struct createAccountParams di sqlc/account.sql.go
	Owner string `json:"owner" binding:"required"`
	// terapkan custom validate yg tlh kita buat dg memanggil tagnya (api/server.go)
	Currency string `json:"currency" binding:"required,currency"`
	/* ==input param ini didptkan dari body HTTP request berupa json makanya ada json:
	lalu binding untuk validasi memakai library go-playground/validator/v10
	== currency hanya USD & EUR, cara validasinya lihat https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags
	bagian One Of itu u/ cek value input jd validasi USD / EUR
	, oneof=USD EUR dlm binding

	*/
}

// func dg server pointer receiver. paramnya objek gin.Context as input
// karena (see) POST() di part param ada HandlerFunc dideklarasikan as func dg Context input
// umumnya when use gin everything we do in side handler will melibatkan obj Context ini, juga u/ read input param & write out responses
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	// ShouldBindingJSON get data from json body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// jika ada error mean client invalid data, send 400 (bad req) to client
		// param 1 status kode (400), param 2 JSON obj send to client (send error dg obj key value to client) selain itu func param 2 will not just account handler. errorResponse dibuat di api/server.go
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Jika input data valid
	// insert new account to db
	arg := db.CreateAccountParams{
		// owner didpt dr request dg nama Owner
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		// convert error pq
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		// internal issue when try insert to db. ke client 500, & error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// kembalikan semuanya
		return
	}
	// if no error, akun berhsl dibuat. kirim status 200 & objek dr account yg dibuat
	ctx.JSON(http.StatusOK, account)
}

//define struct getAccountRequest to store input param
type getAccountRequest struct {
	// krn id is url param for get it use  "uri" (part of Bind Uri at gin-gonic installation) dan jlnkan ShouldBindUri() for bind all URL param into struct
	// validasi id !< 0 so use min
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	// ShouldBindingUri get data from url param
	if err := ctx.ShouldBindUri(&req); err != nil {
		// jika ada error mean client invalid data, send 400 (bad req) to client
		// param 1 status kode (400), param 2 JSON obj send to client (send error dg obj key value to client) selain itu func param 2 will not just account handler. errorResponse dibuat di api/server.go
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// jika tdk error
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// error maka ada 2 kondisi.
		// 1) jika id nya tdk ditemukan
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			// kembalikan semuanya
			return
		}
		// 2) internal error saat query db
		// internal issue when try insert to db. ke client 500, & error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// kembalikan semuanya
		return
	}
	// account = db.Account()//cek response bodynya kosong hrs failed
	// if tdk ada error dan id ditemukan
	ctx.JSON(http.StatusOK, account)
}

// LIST ACCOUNT WITH PAGINATION (QUERY STRING)
//define struct listAccountRequest to store input param
type listAccountRequest struct {
	// krn id is query string for get it use  "form" (part of Only Bind Query String at gin-gonic installation) dan jlnkan ShouldBindQuery() for bind all query string into struct
	// validasi id !< 0 so use min
	PageID int32 `form:"page_id" binding:"required,min=1"`
	// ukuran dari 1 page mau berapa data (kasus ini sdktnya 5 plng banyak 10 record)
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	// ShouldBindingQuery u/ get data from query string
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// jika ada error mean client invalid data, send 400 (bad req) to client
		// param 1 status kode (400), param 2 JSON obj send to client (send error dg obj key value to client) selain itu func param 2 will not just account handler. errorResponse dibuat di api/server.go
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//
	arg := db.ListAccountsParams{
		// 2 field yaitu limit & offset
		Limit: req.PageSize,
		// offset is jumlah record dr db yg hrs di skip, hitung dr page id & page size
		Offset: (req.PageID - 1) * req.PageSize,
	}
	// jika tdk error
	accounts, err := server.store.ListAccounts(ctx, arg) // ListAccounts param 2 bth ListAccountsParams as input
	if err != nil {
		// internal issue. ke client 500, & error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// kembalikan semuanya
		return
	}
	// if tdk ada error dan id ditemukan
	ctx.JSON(http.StatusOK, accounts)
}
