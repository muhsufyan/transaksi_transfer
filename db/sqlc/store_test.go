package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// mentransfer
func TestTransferTx(t *testing.T) {
	store := NewStore((testDB))
	// transfer dari akun1 ke akun2
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// debug
	fmt.Println(">> sblm transaksi(transfer) : ", account1.Balance, account2.Balance)

	// transaksi hrs hati" karena hrs handle concurrency, salah satu solusinya run dg concurrent go routine(di bhsa lain disbt multi thread/asynchronous)
	// run concurrent transfer transaction
	n := 5              //run 5 concurrent transfer transaction
	amount := int64(10) //1 kali transfer sebsr 10 (e.g. $10)
	// now testify tdk bisa cek karena fungsi run go routine berbeda eksekusinya dg TestTransferTx. jd untuk verify/cek error dan result is send them to the main go routine yaitu dimana tmpt test kita run
	// untuk melakukan itu kita bth channel (untuk menghubkan concurrent go routine) so share data aman tanpa terjd locking
	errs := make(chan error)               //BUAT CHANNEL error. untuk receive error from difference go routine to main routine.
	results := make(chan TransferTxResult) //BUAT CHANNEL hsl transfer. untuk receive hsl transfer

	for i := 0; i < n; i++ {
		// start routine.
		go func() {
			// lakukan dan simpan transaksi
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				// transfer amount uang dari akun1 ke akun2
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err       // send error to the errors channel
			results <- result // send result to the results channel
		}() //agar jln hrs ditambah tanda ()
	}
	// cek result, mengecek transfer & entry obj yg tlh dibuat
	for i := 0; i < n; i++ {
		err := <-errs //receive error from channel errors
		require.NoError(t, err)
		result := <-results //receive result from channel results
		require.NotEmpty(t, result)
		// cek transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// cek transfer record tlh tersimpan (dibuat) di db
		_, err = store.GetTransfer(context.Background(), transfer.ID) //karena obj Queries embedded inside Store so GetTransfer juga ada di Store(sprti inheritance)
		require.NoError(t, err)

		// cek akun entries of the result(cek entries). entry pengirim
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		// get akun entry from db untuk make sure that data bnr" sdh dibuat
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// entry penerima
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		// get akun entry from db untuk make sure that data bnr" sdh dibuat
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
	}
}
