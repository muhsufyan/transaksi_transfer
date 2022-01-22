package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/muhsufyan/transaksi_transfer/util"

	"github.com/stretchr/testify/require"
)

// membuat data scra random untuk testing
func createRandomAccount(t *testing.T) Account {
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
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	// buat data scra random & simpan ke db
	account1 := createRandomAccount(t)
	// cari data akun1 dari db yg tlh disimpan ke db
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	// testing
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	// cek semua data di variabel account2 sehrsnya == data di variabel account1. pengecekkan dilakukan satu persatu
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	// u/ field timestep sprti CreatedAt, yg dicari adlh perbedaan waktunya di kasus ini dlm bntk detik
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	// run TestGetAccount
}

func TestUpdateAccount(t *testing.T) {
	// buat data scra random & simpan ke db
	account1 := createRandomAccount(t)
	// untuk menangkap data yg diupdate
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: account1.Balance,
	}
	// jalankan query update ke db
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	// cek semua data di variabel account2 sehrsnya == data di variabel account1. pengecekkan dilakukan satu persatu
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	// karena yg diupdate hanya balance
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	// u/ field timestep sprti CreatedAt, yg dicari adlh perbedaan waktunya di kasus ini dlm bntk detik
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	// run TestUpdateAccount
}

func TestDeleteAccount(t *testing.T) {
	// buat data scra random & simpan ke db
	account1 := createRandomAccount(t)
	// jlnkan query delete ke db
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	// get data dari id yg telah dihapus
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	// var err & error seharusnya adlh tdk ada data tsb
	require.EqualError(t, err, sql.ErrNoRows.Error())
	// data yg dicari hrsnya tdk ada karena sdh dihapus
	require.Empty(t, account2)
	// run TestDeleteAccount
}

func TestListAccount(t *testing.T) {
	// buat 10 account baru dg data yg digenerate random & simpan ke db
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	// tangkap 5 data dari 10
	arg := ListAccountsParams{
		// limit & offset = 5 artinya skip 5 record pertama dan return 5 record selanjutnya
		Limit:  5,
		Offset: 5,
	}
	// dr 10 yg dibuat get 5 record saja yg didpt dr db
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	// returned akun 5
	require.Len(t, accounts, 5)
	// iterate list of account
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
	//  run TestListAccount
}

// run package test
