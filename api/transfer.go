package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
)

// mengecek (membandingkan) mata uang/currency pentransfer dan penerima (mata uangnya hrs sama)
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	// get akun from db
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		// skenario 1 jika akun tdk ada didb
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		// skenario 2 jika unexpected error occur
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	// jika tdk ada error tp currencynya tdk sama
	if account.Currency != currency {
		err := fmt.Errorf("mata uang dari akun [%d] tidak sama : %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	// jika tdk ada mslh & valid
	return true
}

// for store the create account request
type transferRequest struct {
	// money going out
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	// money going in
	ToAccountID int64 `json:"to_account_id" binding:"required,min=1"`
	// jumlah uang yg ditransfer. uang ditransfer hrs > 0
	Amount   int64  `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR RP"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	// ShouldBindingJSON get data from json body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// jika ada error mean client invalid data, send 400 (bad req) to client
		// param 1 status kode (400), param 2 JSON obj send to client (send error dg obj key value to client) selain itu func param 2 will not just account handler. errorResponse dibuat di api/server.go
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// validasi mata uang nya sama (USD => USD) untuk pentransfer
	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	// validasi mata uang nya sama (USD => USD) untuk penerima
	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}
	// Jika input data valid. insert new transfer to db
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	// melakukan transaksi transfer uang
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		// internal issue when try insert to db. ke client 500, & error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// kembalikan semuanya
		return
	}
	// if no error, berhsl melakukan transfer. kirim status 200 & objek dr transfre yg dibuat
	ctx.JSON(http.StatusOK, result)
}
