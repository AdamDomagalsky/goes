package token

import (
	"testing"
	"time"

	"github.com/AdamDomagalsky/goes/bank/util"
	"github.com/stretchr/testify/require"
)

func TestPASETOMaker(t *testing.T) {
	maker, err := NewPASETOMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPASETOToken(t *testing.T) {
	maker, err := NewPASETOMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestExpiredPASETOTokens(t *testing.T) {
	_, err := NewPASETOMaker(util.RandomString(64))
	require.ErrorContains(t, err, "invalid key size: must be exactly 32 characters long")
}
