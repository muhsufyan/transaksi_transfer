package db

import (
	"context"
	"testing"
	"time"

	"github.com/muhsufyan/transaksi_transfer/util"
	"github.com/stretchr/testify/require"
)

// membuat data scra random untuk testing
func createRandomUser(t *testing.T) User {
	// hash pass, generate new pass (6 karakter)
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	// data random untuk testing
	arg := CreateUserParams{
		// kita isi datanya
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	// buat user baru dan simpan ke db(melalui CreateUser()) dg data dari data random yg tlh dibuat
	user, err := testQueries.CreateUser(context.Background(), arg)
	// testing
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	// when user pertama kali dibuat nilai default timestamp passwordnya shrsnya 0 (so use IsZero())
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// get user
func TestGetUser(t *testing.T) {
	// buat data random & simpan ke db
	user1 := createRandomUser(t)
	// get data yg baru dibuat td di db berdsrkan username
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	// WithinDuration untuk compare nilai (kasus ini timestamp) dg perbedaan yg kecil (kasus ini dlm satuan detik)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
