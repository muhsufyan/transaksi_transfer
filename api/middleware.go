package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muhsufyan/transaksi_transfer/token"
)

const (
	authorizationHeaderKey = "otorisasi"
	// misal hanya support 1 type otorisasi yaitu bearer token
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// func ini sbnrnya bkn middleware tp heigher-order func. param is token maker interface
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	// return func otentifikasi middleware
	return func(ctx *gin.Context) {
		// extract otorisasi header dr request. param otorisasi header key
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		// empty. client tdk punya header ini
		if len(authorizationHeader) == 0 {
			err := errors.New("header otorisasi tidak ada")
			// abort request
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// ada header otorisasi
		// split header otorisasi dg spasi
		fields := strings.Fields(authorizationHeader)
		// shrsnya dihsl sdktnya 2 element
		if len(fields) < 2 {
			err := errors.New("format headernya invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// elemen pertama otorisasi dlh field slice
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("tipe otorisasi tdk didukung %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// tipe otorisasinya sama
		accessToken := fields[1]
		// parse dan verify access token to get payload(data)
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// token is valid
		// store payload in context before passing payload to next handler
		// store payload in context
		ctx.Set(authorizationPayloadKey, payload)
		// forward request to next handler
		ctx.Next()
	}
}
