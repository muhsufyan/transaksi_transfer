package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// pesan token yg error
var (
	ErrInvalidToken = errors.New("Token is invalid")
	ErrExpiredToken = errors.New("Token has expired")
)

// berisi data dr token
type Payload struct {
	// UUID untuk token yg dibuat karena setiap token hrs punya unik key
	ID uuid.UUID `json:"id"`
	// mengidentifikasi pemilik dr token
	Username string `json:"username"`
	// waktu token dibuat(waktu pembuatan token)
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// membuat token baru dg payload (data) berupa username dan durasi(durasi aktifnya suatu token)
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	// generate unik token id
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	// tdk ada error
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, err
}

// cek payloadnya valid/tdk. lihat token/maker.go jwt.NewWithClaims param 2 payload jd tdk error karena param 2 perlu fungsi ini
func (payload *Payload) Valid() error {
	// cek tokennya expired
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
