package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/adifahmi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomUserParams() CreateUserParams {
	hashedPass, _ := util.HashPassword(util.RandomString(6, false))
	return CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
}

func createRandomUser(arg CreateUserParams) (sql.Result, error) {
	res, err := testQueries.CreateUser(context.Background(), arg)
	return res, err
}

func getSingleUserByUsername(username string) (User, error) {
	user, err := testQueries.GetUserByUsername(context.Background(), username)
	return user, err
}

func getSingleUserByID(userID int64) (User, error) {
	user, err := testQueries.GetUserByID(context.Background(), userID)
	return user, err
}

func createAndGetUser() (User, error) {
	arg := randomUserParams()
	res, _ := createRandomUser(arg)
	userId, _ := res.LastInsertId()
	return getSingleUserByID(userId)
}

func TestCreateUser(t *testing.T) {
	arg := randomUserParams()

	res, err := createRandomUser(arg)
	require.NoError(t, err)
	require.NotEmpty(t, res)
}

func TestGetUser(t *testing.T) {
	arg := randomUserParams()
	res, _ := createRandomUser(arg)
	userID, _ := res.LastInsertId()
	user, err := getSingleUserByID(userID)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	t.Log("Username:", user.Username)
	require.Equal(t, userID, user.ID)

	user2, err2 := getSingleUserByUsername(arg.Username)
	require.NoError(t, err2)
	require.Equal(t, user.Username, user2.Username)
}
