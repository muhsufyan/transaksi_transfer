package token

import "time"

// general token maker interface untuk manage pembuatan dan verifikasi token
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
