package db

import (
	"context"
	"database/sql"
	"fmt"
)
// ganti obj struct Store dg interface
type Store interface{
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}
// ganti Store jd SQLStore.SQLStore menyediakan semua fungsi to eksekusi SQL queries (postgres) dan transaksi
// berisi semua fungsi query db (sblmnya query hanya dpt dijlnkan untuk 1 tabel saja, dg store itu ibaratnya variabel global untuk query(dlm konsep state management di reactjs)) dan transaksi
type SQLStore struct {
	// struct Queries tdk mendukung query transaksi jd solusinya dg Store ini, caranya dg embedding kedlm struct store (extend atau kita sbt inheritance di go disbt composition)
	*Queries
	// query db transaction
	db *sql.DB
}
// return pointer (*Store) diubah jd return interface (Store)
// NewStore create a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execute generic db transaction dg param context & fungsi callback as input. execTx akan mengeksekusi fungsi transaksi db
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	//obj Queries untuk transaksi. param 2 bisa diisi &sql.TxOptions() untuk custom level transaksi tp sekarang we tdk perlu (use default level isolasi). tx adlh transfer transaksi
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	// run query transaksi
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("txt err : %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// menangkap semua input param dari transaksi transfer
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// hsl dr transaksi transfer
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_transfer"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// melakukan transaksi transfer uang dr 1 akun ke akun lain
// yg dilakukan is membuat record transfer, add account entries, dan update akun balance dg 1 transaksi db
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// create & run new db transaksi
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// query untuk transaksi membuat transfer, param 2 tangkap data yg dimasukkan
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, //uang ditransfer /keluar
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// update akun (transaksi) yg pertama dilakukan adlh yg id nya paling kecil
		// si pentransfer(from akun) di update dulu baru si penerima(to akun)
		if arg.FromAccountID < arg.ToAccountID {
			// update from akun(pengirim transfer) dulu baru to akun(menerima transfer)
			// -arg.Amount karena dia mentransfer uang, sedangkan arg.Amount karena dia menerima uang
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else { // si penerima(to akun) di update dulu baru si pentransfer(from akun)
			// atau toAccount diupdate sblm fromAccount
			// arg.Amount karena dia menerima uang, sedangkan -arg.Amount karena dia mentransfer uang
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

		}

		return nil
	})
	return result, err
}

// REFACTOR
// add money to 2 account
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	//param untuk return akunya stlh diupdate
	// add amount1 to account1 balance. // update akun
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	// add amount2 to account2 balance
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
