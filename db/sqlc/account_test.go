package db

import (
	"context"
	"testing"

	"github.com/muhsufyan/transaksi_transfer/util"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	// data random untuk testing
	arg := CreateAccountParams{
		// kita isi datanya
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	// buat akun baru dan simpan ke db(melalui CreateAccount()) dg data dari data random yg tlh dibuat
	account, err := testQueries.CreateAccount(context.Background(), arg)
	// testing
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
