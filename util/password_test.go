package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

/*
unit test, make sure fungsi HashPassword dan CheckPassword work sprti ekspektasi kita (ke duanya ada di util/password.go)
*/
func TestPassword(t *testing.T) {
	// generate new pass dg 6 karakter
	password := RandomString(6)
	// hash pass nya
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	// cek pass dg pass yg tlh dihash hrs sama
	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)
	// test case saat incorrect pass
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
