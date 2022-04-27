package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error)
}

type SqlStore struct {
	db *sql.DB
	*Queries
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

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

		// Actual transfer here
		if arg.FromAccountID < arg.ToAccountID {
			// deduct first then add
			result.FromAccount, result.ToAccount, err = updateBothBalance(ctx, q, arg.FromAccountID, -arg.Ammount, arg.ToAccountID, arg.Ammount)
		} else {
			// add first then deduct
			result.ToAccount, result.FromAccount, err = updateBothBalance(ctx, q, arg.ToAccountID, arg.Ammount, arg.FromAccountID, -arg.Ammount)
		}

		return err
	})

	return result, err
}

// updateBothBalance updates both accounts balance (sender and receiver)
// sender could be account1 or account2 and vice versa
func updateBothBalance(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	ammount1 int64,
	accountID2 int64,
	ammount2 int64,
) (account1 Account, account2 Account, err error) {
	err = q.UpdateAccountBalance(context.Background(), UpdateAccountBalanceParams{
		ID:      accountID1,
		Ammount: ammount1,
	})
	if err != nil {
		return
	}
	account1, err = q.GetAccount(context.Background(), accountID1)
	if err != nil {
		return
	}

	err = q.UpdateAccountBalance(context.Background(), UpdateAccountBalanceParams{
		ID:      accountID2,
		Ammount: ammount2,
	})
	if err != nil {
		return
	}
	account2, err = q.GetAccount(context.Background(), accountID2)
	if err != nil {
		return
	}
	return
}
