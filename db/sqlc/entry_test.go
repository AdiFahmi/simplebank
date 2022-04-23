package db

import (
	"context"
	"database/sql"
	"testing"

	"example.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(arg CreateEntryParams) (sql.Result, error) {
	res, err := testQueries.CreateEntry(context.Background(), arg)
	return res, err
}

func getSingleEntry(entryID int64) (Entry, error) {
	entry, err := testQueries.GetEntry(context.Background(), entryID)
	return entry, err
}

func TestCreateRandomEntry(t *testing.T) {
	account, _ := createAndGetAccount()
	arg := CreateEntryParams{
		AccountID: account.ID,
		Ammount:   util.RandomMoney(),
	}

	res, err := createRandomEntry(arg)
	require.NoError(t, err)
	require.NotEmpty(t, res)
}

func TestGetEntry(t *testing.T) {
	account, _ := createAndGetAccount()
	arg := CreateEntryParams{
		AccountID: account.ID,
		Ammount:   util.RandomMoney(),
	}
	res, _ := createRandomEntry(arg)
	entryID, _ := getIdFromResult(res)
	entry, err := getSingleEntry(entryID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entryID, entry.ID)
	require.IsType(t, int64(0), entry.Ammount)
}

func TestListEntries(t *testing.T) {
	account, _ := createAndGetAccount()
	for i := 0; i <= 10; i++ {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Ammount:   util.RandomMoney(),
		}
		_, _ = createRandomEntry(arg)
	}

	arg2 := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
