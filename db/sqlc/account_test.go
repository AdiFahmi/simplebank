package db

import (
	"context"
	"database/sql"
	"testing"

	"example.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomParams() CreateAccountParams {
	return CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func createRandomAccount(arg CreateAccountParams) (sql.Result, error) {
	res, err := testQueries.CreateAccount(context.Background(), arg)
	return res, err
}

func getSingleAccount(accID int64) (Account, error) {
	acc, err := testQueries.GetAccount(context.Background(), accID)
	return acc, err
}

func createAndGetAccount() (Account, error) {
	arg := randomParams()
	res, _ := createRandomAccount(arg)
	accID, _ := res.LastInsertId()
	return getSingleAccount(accID)
}

func TestCreateAccount(t *testing.T) {
	arg := randomParams()

	res, err := createRandomAccount(arg)
	require.NoError(t, err)
	require.NotEmpty(t, res)
}

func TestGetAccount(t *testing.T) {
	arg := randomParams()
	res, _ := createRandomAccount(arg)
	accID, _ := res.LastInsertId()
	acc, err := getSingleAccount(accID)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	t.Log("Acc owner:", acc.Owner)
	require.Equal(t, accID, acc.ID)
	require.IsType(t, int64(0), acc.Balance)
}

func TestUpdateAccount(t *testing.T) {
	arg := randomParams()
	res, _ := createRandomAccount(arg)
	accID, _ := res.LastInsertId()

	updateArg := UpdateAccountParams{
		ID:      accID,
		Balance: util.RandomMoney(),
	}

	err := testQueries.UpdateAccount(context.Background(), updateArg)
	require.NoError(t, err)

	updatedAcc, err := getSingleAccount(accID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	require.Equal(t, updateArg.Balance, updatedAcc.Balance)
}

func TestDeleteAccount(t *testing.T) {
	arg := randomParams()
	res, _ := createRandomAccount(arg)
	accID, _ := res.LastInsertId()

	err := testQueries.DeleteAccount(context.Background(), accID)
	require.NoError(t, err)

	acc, err := getSingleAccount(accID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		arg := randomParams()
		_, _ = createRandomAccount(arg)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		t.Log("Acc owner:", account.Owner)
	}
}

func TestAddAccountBalance(t *testing.T) {
	account, _ := createAndGetAccount()
	t.Log("Init balance:", account.Balance)
	addAmmount := util.RandomMoney()
	t.Log("added Ammount:", addAmmount)

	err := testQueries.UpdateAccountBalance(context.Background(), UpdateAccountBalanceParams{
		ID:      account.ID,
		Ammount: addAmmount,
	})
	require.NoError(t, err)

	updatedAcc, err := getSingleAccount(account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	t.Log("Updated balance:", updatedAcc.Balance)
	require.Equal(t, account.Balance+addAmmount, updatedAcc.Balance)
}

func TestDeductAccountBalance(t *testing.T) {
	account, _ := createAndGetAccount()
	t.Log("Init balance:", account.Balance)
	deductAmmount := util.RandomMoney()
	t.Log("deduct Ammount:", deductAmmount)

	err := testQueries.UpdateAccountBalance(context.Background(), UpdateAccountBalanceParams{
		ID:      account.ID,
		Ammount: -deductAmmount,
	})
	require.NoError(t, err)

	updatedAcc, err := getSingleAccount(account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	t.Log("Updated balance:", updatedAcc.Balance)
	require.Equal(t, account.Balance-deductAmmount, updatedAcc.Balance)
}
