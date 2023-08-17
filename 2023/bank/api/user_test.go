package api

import (
	"context"
	"fmt"
	"testing"

	db "github.com/AdamDomagalsky/goes/2023/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/2023/bank/util"
	"github.com/stretchr/testify/require"
)

func CreateUser(store db.Store) (db.User, error) {
	return store.CreateUser(context.Background(), db.CreateUserParams{})
}

func randomUser(t *testing.T, store db.Store) (user db.User, password string) {
	password = util.RandomString(12)
	hash, err := util.HashPassword(password)
	require.NoError(t, err)

	user, err = createUser(store, db.CreateUserParams{
		Username: util.RandomOwner(),
		Password: hash,
		Email:    util.RandomEmail(),
		FullName: fmt.Sprintf("%s %s", util.RandomString(3), util.RandomString(8)),
	})
	require.NoError(t, err)

	return user, password

}

func createUser(store db.Store, params db.CreateUserParams) (db.User, error) {
	return store.CreateUser(context.Background(), params)
}
