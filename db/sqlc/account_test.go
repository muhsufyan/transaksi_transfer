package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	// data dummy untuk testing
	arg := CreateAccountParams{
		// kita isi datanya
		Owner:    "owner testing 1",
		Balance:  100,
		Currency: "USD",
	}
	// buat akun baru dan simpan ke db(melalui CreateAccount()) dg data dari data dummy yg tlh dibuat
	account, err := testQueries.CreateAccount(context.Background(), arg)
	// langkah selanjutnya cek hsl test dg apa yg kita inginkan. untuk itu kita perlu install testify
	// require.NoError paramnya (passing) objek t dan error. ini akan mengecek errornya (actually) itu nil dan akan otomatis fail jika tdk nil(ada error). ini mengetes return error
	require.NoError(t, err)
	// returnya objek shrsnya tdk kosong, objek account shrsnya tdk kosong
	require.NotEmpty(t, account)
	// cek apakah yg diinputkan itu sama dg apa yg diharapkan(actually)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	// cek apakah Id nya dibuat otomatis (autoincrement)
	require.NotZero(t, account.ID)
	// cek apakah CreateAt adalah waktu saat ini dan berupa timestep
	require.NotZero(t, account.CreatedAt)
}
