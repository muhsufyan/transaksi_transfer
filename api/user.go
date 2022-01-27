package api

import (
	"database/sql"
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
type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// convert obj input db.User jd userResponse
func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	// ShouldBindingJSON get data from json body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// jika ada error mean client invalid data, send 400 (bad req) to client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// hashedPassword, err := util.HashPassword("xyz") //unacceptable, hrsnya saat test adlh failed
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
	// arg = db.CreateUserParams{} //datanya kosong sehrsnya failed when testing
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
	dataResponse := newUserResponse(user)
	// if no error, user berhsl dibuat. kirim status 200 & objek dr user yg dibuat
	ctx.JSON(http.StatusOK, dataResponse)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	// token ini dibuat oleh token maker interface
	AccessToken string `json:"access_token"`
	// informasi logged user
	User userResponse `json:"user"`
}

// api handler untuk login
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	// paramnya pointer ke obj req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// jika ada error mean client invalid data, send 400 (bad req) to client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// tdk ada error find data user di db
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		// ada 2 possible cases
		// 1) username tdk ada di db
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// 2) unexpected error occurs when talking to db
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// jika good kita cek pass usernya bnr/tdk ?
	err = util.CheckPassword(req.Password, user.HashedPassword)
	// passwordnya salah
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	// pass-nya bnr maka buat/generate token
	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	// jika terjd unexpected error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// jika good build obj loginUserResponse
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
