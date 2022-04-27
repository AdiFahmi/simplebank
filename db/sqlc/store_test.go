package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1, _ := createAndGetAccount()
	account2, _ := createAndGetAccount()
	t.Log("Init balance", account1.Balance, account2.Balance)

	totalLoop := 1
	ammount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < totalLoop; i++ {
		txName := fmt.Sprintf("tx-%d", i)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Ammount:       ammount,
			})

			errs <- err
			results <- result
		}()
	}

	// check error and result
	for i := 0; i < totalLoop; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, ammount, transfer.Ammount)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -ammount, fromEntry.Ammount)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, ammount, toEntry.Ammount)
		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		t.Log("tx Balance", account1.Balance, account2.Balance)
		diff1 := account1.Balance - fromAccount.Balance // money from acc1
		diff2 := toAccount.Balance - account2.Balance   // money into acc2
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%ammount == 0)

		k := int(diff1 / ammount)
		require.True(t, k >= 1 && k <= totalLoop)
	}

	// check final balance
	updatedAccount1, _ := store.GetAccount(context.Background(), account1.ID)
	updatedAccount2, _ := store.GetAccount(context.Background(), account2.ID)
	t.Log("Final balance", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(totalLoop)*ammount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(totalLoop)*ammount, updatedAccount2.Balance)
}
