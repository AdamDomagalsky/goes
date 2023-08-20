package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AdamDomagalsky/goes/2023/bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	user, _ := CreateRandomUser(t, testEnv.Store)
	CreateRandomAccount(t, testEnv.Store, user)
}
func TestGetAccount(t *testing.T) {
	user, _ := CreateRandomUser(t, testEnv.Store)
	account1 := CreateRandomAccount(t, testEnv.Store, user)
	account2, err := testEnv.Store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
func TestUpdateAccount(t *testing.T) {
	user, _ := CreateRandomUser(t, testEnv.Store)
	account1 := CreateRandomAccount(t, testEnv.Store, user)
	newBalance := util.RandomMoney()
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: newBalance,
	}

	account2, err := testEnv.Store.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, newBalance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
func TestDeleteAccount(t *testing.T) {
	user, _ := CreateRandomUser(t, testEnv.Store)
	account1 := CreateRandomAccount(t, testEnv.Store, user)
	err := testEnv.Store.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account2, err := testEnv.Store.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		user, _ := CreateRandomUser(t, testEnv.Store)
		lastAccount = CreateRandomAccount(t, testEnv.Store, user)
	}
	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testEnv.Store.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 1)
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestListAccountsNegativeOffset(t *testing.T) {
	arg := ListAccountsParams{
		Offset: -35,
	}
	accounts, err := testEnv.Store.ListAccounts(context.Background(), arg)
	require.EqualError(t, err, "pq: OFFSET must not be negative")
	require.Len(t, accounts, 0)
}
