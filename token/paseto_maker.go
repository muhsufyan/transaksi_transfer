package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// paseto token maker. implement token maker interface
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// buat new PasetoMaker.
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// cek jumlah karakter simetri key tdk sama dg jumlah karakter key dr algo chacha poly yg used u/ enkripsi payloadnya
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	// jika sama, buat obj PasetoMaker yg baru.
	// symmetricKey converted to byte slice
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// 2 fungsi brkt implement dr interface Maker digunakan agar return pd NewPasetoMaker
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// create new token payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// param 3 itu optional footer we not need so nil
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)

}
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	// store decrypted data
	payload := &Payload{}

	// param 3 itu optional footer we not need so nil
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
