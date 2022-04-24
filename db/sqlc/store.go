package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error)
}

type SqlStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *SqlStore) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		txName := ctx.Value(txKey)

		// create transfer log
		fmt.Println(txName, "create transfer log")
		createTfRes, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Ammount:       arg.Ammount,
		})
		if err != nil {
			return err
		}
		transferID, err := createTfRes.LastInsertId()
		if err != nil {
			return err
		}
		// append transfer log to result
		result.Transfer, err = q.GetTransfer(ctx, transferID)
		if err != nil {
			return err
		}

		// create entry FROM log
		fmt.Println(txName, "create entry FROM log")
		createEntryFromRes, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Ammount:   -arg.Ammount,
		})
		if err != nil {
			return err
		}
		entryFromID, err := createEntryFromRes.LastInsertId()
		if err != nil {
			return err
		}
		result.FromEntry, err = q.GetEntry(ctx, entryFromID)
		if err != nil {
			return err
		}

		// create entry TO log
		fmt.Println(txName, "create entry TO log")
		createEntryToRes, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Ammount:   arg.Ammount,
		})
		if err != nil {
			return err
		}
		entryToID, err := createEntryToRes.LastInsertId()
		if err != nil {
			return err
		}
		result.ToEntry, err = q.GetEntry(ctx, entryToID)
		if err != nil {
			return err
		}

		// Update account balance
		// Deduct from account1
		fmt.Println(txName, "update account1 balance")
		err = q.UpdateAccountBalance(context.Background(), UpdateAccountBalanceParams{
			ID:      arg.FromAccountID,
			Ammount: -arg.Ammount,
		})
		if err != nil {
			return err
		}
		result.FromAccount, err = q.GetAccount(context.Background(), arg.FromAccountID)
		if err != nil {
			return err
		}

		// Add to account2
		fmt.Println(txName, "update account2 balance")
		err = q.UpdateAccountBalance(context.Background(), UpdateAccountBalanceParams{
			ID:      arg.ToAccountID,
			Ammount: arg.Ammount,
		})
		if err != nil {
			return err
		}
		result.ToAccount, err = q.GetAccount(context.Background(), arg.ToAccountID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
