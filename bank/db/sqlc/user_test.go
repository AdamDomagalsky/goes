package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t, testEnv.Store)
}

func TestGetUser(t *testing.T) {
	user1, _ := CreateRandomUser(t, testEnv.Store)
	user2, err := testEnv.Store.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserPartial(t *testing.T) {
	store := testEnv.Store
	b4User, _ := CreateRandomUser(t, testEnv.Store)
	newEmail := "elo@o2.pl"
	afterUser, err := store.UpdateUser(context.Background(), UpdateUserParams{
		Email:    sql.NullString{newEmail, true},
		Username: b4User.Username,
	})
	require.NoError(t, err)

	require.Equal(t, b4User.Username, afterUser.Username)
	require.Equal(t, b4User.FullName, afterUser.FullName)
	require.Equal(t, b4User.Password, afterUser.Password)
	require.Equal(t, b4User.PasswordChangedAt, afterUser.PasswordChangedAt)
	require.Equal(t, b4User.CreatedAt, afterUser.CreatedAt)

	require.NotEqual(t, b4User.Email, afterUser.Email)
	require.Equal(t, newEmail, afterUser.Email)
}

func TestUpdateUserPartialExample(t *testing.T) {
	store := testEnv.Store
	b4User, _ := CreateRandomUser(t, testEnv.Store)
	newFullname := "Full Name"
	afterUser, err := store.UpdateUserCaseExample(context.Background(), UpdateUserCaseExampleParams{
		FullName:    newFullname,
		SetFullName: true,
		Username:    b4User.Username,
	})
	require.NoError(t, err)

	require.Equal(t, b4User.Username, afterUser.Username)
	require.Equal(t, b4User.Password, afterUser.Password)
	require.Equal(t, b4User.PasswordChangedAt, afterUser.PasswordChangedAt)
	require.Equal(t, b4User.CreatedAt, afterUser.CreatedAt)
	require.Equal(t, b4User.Email, afterUser.Email)

	require.NotEqual(t, b4User.FullName, afterUser.FullName)
	require.Equal(t, newFullname, afterUser.FullName)
}
