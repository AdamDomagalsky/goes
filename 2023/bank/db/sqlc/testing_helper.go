package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/AdamDomagalsky/goes/2023/bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T, store Store) Account {
	user, _ := CreateRandomUser(t, store)
	result, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)

	return result
}

func CreateRandomUser(t *testing.T, store Store) (user User, password string) {
	password = util.RandomString(12)
	hash, err := util.HashPassword(password)
	require.NoError(t, err)

	user, err = store.CreateUser(context.Background(), CreateUserParams{
		Username: util.RandomOwner(),
		Password: hash,
		Email:    util.RandomEmail(),
		FullName: fmt.Sprintf("%s %s", util.RandomString(3), util.RandomString(8)),
	})
	require.NoError(t, err)

	return user, password

}
