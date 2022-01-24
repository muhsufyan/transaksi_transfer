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

	// run concurrent transfer transaction. create 5 go routine to execute 5 concurrent transfer transaction
	n := 5                                 //run 5 concurrent transfer transaction
	amount := int64(10)                    //1 kali transfer sebsr 10 (e.g. $10)
	errs := make(chan error)               //BUAT CHANNEL error. untuk receive error from difference go routine to main routine.
	results := make(chan TransferTxResult) //BUAT CHANNEL hsl transfer. untuk receive hsl transfer

	for i := 0; i < n; i++ {
		// start routine.
		go func() {
			ctx := context.Background()
			// lakukan dan simpan transaksi
			result, err := store.TransferTx(ctx, TransferTxParams{
				// transfer amount uang dari akun1 ke akun2
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err       // send error to the errors channel
			results <- result // send result to the results channel
		}() //agar jln hrs ditambah tanda ()
	}
	existed := make(map[int]bool) //array yg akan menyimpan nilai berapakali pembagian
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

		// cek akun. cek output account
		fromAccount := result.FromAccount // uang yg keluar(ditransfer)
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// uang pergi kemana(tujuan transfer)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// cek akun balance
		fmt.Println(">> tx, balance setiap stlh transaksi : ", fromAccount.Balance, toAccount.Balance) //hsl balance stlh tiap transaksi
		diff1 := account1.Balance - fromAccount.Balance                                                //jumlah uang going  out to account1
		diff2 := toAccount.Balance - account2.Balance                                                  //jumlah uang going  in to account2
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //amount balance dr akun1 setiap 1 kali transaksi akan berkurang, amount * 2 transaksi, amount * 3 transaksi. misal amount = 10

		// test jumlah transaksi yg dilakukan
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// cek final updated. balance dr akun2
	// get akun1 terupdate dr db
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	// get akun2 terupdate dr db
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> stlh transaksi(transfer) : ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
