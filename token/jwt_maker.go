package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// min karakter token is 32 karakter
const minSecretKeySize = 32

// we memakai algo simetri key
// membuat JWT
type JWTMaker struct {
	secretKey string
}

// membuat JWTMaker yg baru
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size : minimal token adlh %d karakter ", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// 2 fungsi brkt implement dr interface Maker ada di token/maker.go
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// create new token payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// create new jwtToken.param 1 signing method/algoritma yg kita gunakan, param 2 claim adlh payload (data) embedded dlm token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// generate token string. passing secret key stlh converted jd byte slice
	return jwtToken.SignedString([]byte(maker.secretKey))
}
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// algo yg digunakan
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		// algo yg digunakan sama ? (algo encode & decode hrs sama)
		if !ok {
			return nil, ErrInvalidToken
		}
		// algonya sama, secretKey di convert ke byte slice
		return []byte(maker.secretKey), nil
	}
	// parse token
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// ada 2 skenario yaitu token invalid, token expired
		// convert err
		verr, ok := err.(*jwt.ValidationError)
		// errornya expired
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		// errornya invalid
		return nil, ErrInvalidToken
	}
	// semuanya good. get data payload dg convert jd obj payload
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	// jika ok
	return payload, nil
}
