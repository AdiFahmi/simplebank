package db

import (
	"context"
	"database/sql"
	"testing"

	"example.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(arg CreateTransferParams) (sql.Result, error) {
	return testQueries.CreateTransfer(context.Background(), arg)
}

func TestCreateTransfer(t *testing.T) {
	account1, _ := createAndGetAccount()
	account2, _ := createAndGetAccount()

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Ammount:       util.RandomMoney(),
	}

	result, err := createRandomTransfer(arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

func TestGetTransfer(t *testing.T) {
	account1, _ := createAndGetAccount()
	account2, _ := createAndGetAccount()

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Ammount:       util.RandomMoney(),
	}
	result, _ := createRandomTransfer(arg)
	transferID, _ := getIdFromResult(result)

	transfer, err := testQueries.GetTransfer(context.Background(), transferID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transferID, transfer.ID)
	require.IsType(t, int64(0), transfer.Ammount)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, arg.Ammount, transfer.Ammount)
}

func TestListTransfer(t *testing.T) {
	account1, _ := createAndGetAccount()
	account2, _ := createAndGetAccount()

	for i := 0; i < 5; i++ {
		arg := CreateTransferParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Ammount:       util.RandomMoney(),
		}
		_, _ = createRandomTransfer(arg)

		arg2 := CreateTransferParams{
			FromAccountID: account2.ID,
			ToAccountID:   account1.ID,
			Ammount:       util.RandomMoney(),
		}
		_, _ = createRandomTransfer(arg2)
	}

	transferParams := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), transferParams)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.FromAccountID == account2.ID)
	}
}
