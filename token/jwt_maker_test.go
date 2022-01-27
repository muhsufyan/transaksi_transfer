package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/muhsufyan/transaksi_transfer/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	// buat new maker
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	// token valid hanya 1 menit
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	// generate token
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	// make sure token is valid
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	// cek all field of payload obj
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
func TestExpiredJWTToken(t *testing.T) {
	// buat new maker
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	// buat expired token
	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	// verify output token
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// cek kasus invalid token yaitu tdk ada algo yg digunakan
func TestInvalidJWTTokenAlgNone(t *testing.T) {
	// buat new payload dg waktu 1 menit
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	// buat token baru, dg tdk ada algo used, buat payload baru
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	// sign token where tdk memakai algo apapun
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	// buat new maker
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	// verify token yg tlh dibuat sblmnya
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	// output payload hrs nil
	require.Nil(t, payload)
}
