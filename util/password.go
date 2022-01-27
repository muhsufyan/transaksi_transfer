package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// bcrypt hash untuk password
func HashPassword(password string) (string, error) {
	// generate hash dari password. param 1 password of type []byte slice (we convert password dr input jd byte slice), param 2 cost of type int (we use default value yaitu 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// cek passwordnya bnr/salah, param 1 : pass yg dicek param 2: hashed pass to compare
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
