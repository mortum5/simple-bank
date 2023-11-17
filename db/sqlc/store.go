package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	txOptions := &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelReadCommitted,
	}
	tx, err := store.db.BeginTx(ctx, txOptions)

	if err != nil {
		return fmt.Errorf("store: transaction err: %w", err)
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TranferTxResult struct {
	Transfer    Transfer `json:"tranfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

type ValueUpdate func(int64) int64

var TxKey = struct{}{}

func updateBalance(ctx context.Context, q *Queries, userID int64, fn ValueUpdate) (Account, error) {
	account, err := q.GetAccountForUpdate(ctx, userID)
	if err != nil {
		return account, err
	}

	newBalance := fn(account.Balance)
	if newBalance < 0 {
		return account, errors.New("not enough balance")
	}

	account, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      userID,
		Balance: fn(account.Balance),
	})
	if err != nil {
		return account, err
	}

	return account, nil
}

func (store *Store) TranferTx(ctx context.Context, arg TransferTxParams) (TranferTxResult, error) {
	var result TranferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(TxKey)

		dec := func(a int64) int64 {
			return a - arg.Amount
		}

		inc := func(a int64) int64 {
			return a + arg.Amount
		}

		if arg.FromAccountID < arg.ToAccountID {
			fmt.Println(txName, "update from account")
			result.FromAccount, err = updateBalance(ctx, q, arg.FromAccountID, dec)
			if err != nil {
				return err
			}

			fmt.Println(txName, "update to account")
			result.ToAccount, err = updateBalance(ctx, q, arg.ToAccountID, inc)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(txName, "update to account")
			result.ToAccount, err = updateBalance(ctx, q, arg.ToAccountID, inc)
			if err != nil {
				return err
			}

			fmt.Println(txName, "update from account")
			result.FromAccount, err = updateBalance(ctx, q, arg.FromAccountID, dec)
			if err != nil {
				return err
			}
		}

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
