package db

import (
	"context"
	"database/sql"
	"fmt"
)

// berisi semua fungsi query db (sblmnya query hanya dpt dijlnkan untuk 1 tabel saja, dg store itu ibaratnya variabel global untuk query(dlm konsep state management di reactjs)) dan transaksi
type Store struct {
	// struct Queries tdk mendukung query transaksi jd solusinya dg Store ini, caranya dg embedding kedlm struct store (extend atau kita sbt inheritance di go disbt composition)
	*Queries
	// query db transaction
	db *sql.DB
}

// NewStore create a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execute generic db transaction dg param context & fungsi callback as input. execTx akan mengeksekusi fungsi transaksi db
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
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
		// get akun => update balance-nya
		// akun1 as transfer
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}
		// akun2 as ditransfer
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}
