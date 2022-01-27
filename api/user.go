package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
	"github.com/muhsufyan/transaksi_transfer/util"
)

type createUserRequest struct {
	// alphanum mengizinkan spesial karakter
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// data dlm struct ini akan dikirim ke user (jd hashed pass tdk akan dikirim)
type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	// ShouldBindingJSON get data from json body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// jika ada error mean client invalid data, send 400 (bad req) to client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Jika input data valid
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// tdk ada error saat hash password
	// insert new user to db
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	// store to user table
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// convert error pq
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		// internal issue when try insert to db. ke client 500, & error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// kembalikan semuanya
		return
	}
	dataResponse := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	// if no error, user berhsl dibuat. kirim status 200 & objek dr user yg dibuat
	ctx.JSON(http.StatusOK, dataResponse)
}
