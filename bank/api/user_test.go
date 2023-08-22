package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {

	email := util.RandomEmail()
	fullname := util.RandomString(6)
	username := util.RandomString(4)
	password := util.RandomString(6)

	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{
				"email":     email,
				"full_name": fullname,
				"username":  username,
				"password":  password,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				ctx := context.Background()
				testStore.Store.GetAccount(ctx, 3)
				user, err := testStore.Store.GetUser(ctx, username)
				require.NoError(t, err)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "CreateSameUserTwice",
			body: gin.H{
				"email":     email,
				"full_name": fullname,
				"username":  username,
				"password":  password,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				require.Contains(t, recorder.Body.String(), "duplicate key value violates unique constraint")
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case-%d-%s", i, tc.name), func(t *testing.T) {
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			body := bytes.NewReader(data)
			req, err := http.NewRequest(http.MethodPost, "/users", body)
			require.NoError(t, err)
			testServer.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}

}

func TestLoginUserAPI(t *testing.T) {
	user, password := db.CreateRandomUser(t, testServer.store)

	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"username": "imaginaryUser",
				"password": password,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				require.Contains(t, recorder.Body.String(), "no rows in result set")
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"username": user.Username,
				"password": "password",
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), "hashedPassword is not the hash of the given password")
			},
		},
		{
			name: "NoPassword",
			body: gin.H{
				"username": user.Username,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "hashedPassword is not the hash of the given password")
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case-%d-%s", i, tc.name), func(t *testing.T) {
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			body := bytes.NewReader(data)
			req, err := http.NewRequest(http.MethodPost, "/users/login", body)
			require.NoError(t, err)
			testServer.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	decoder := json.NewDecoder(body)
	var bodyUserAccount db.User
	err := decoder.Decode(&bodyUserAccount)
	user.CreatedAt = user.CreatedAt.UTC()
	user.PasswordChangedAt = user.PasswordChangedAt.UTC()
	require.NoError(t, err)

	require.Equal(t, user.Email, bodyUserAccount.Email)
	require.Equal(t, user.FullName, bodyUserAccount.FullName)
	require.Equal(t, user.Username, bodyUserAccount.Username)
	require.Equal(t, user.PasswordChangedAt, bodyUserAccount.PasswordChangedAt)
	require.Equal(t, user.CreatedAt, bodyUserAccount.CreatedAt)
}
