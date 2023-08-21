package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/token"
	"github.com/AdamDomagalsky/goes/bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	user, _ := db.CreateRandomUser(t, testServer.store)
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{
				"owner":    user.Username,
				"currency": util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeKey, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				ctx := context.Background()
				testStore.Store.GetAccount(ctx, 3)
				accounts, err := testStore.Store.ListAccounts(ctx, db.ListAccountsParams{
					Owner:  user.Username,
					Offset: 0,
					Limit:  1,
				})
				require.NoError(t, err)
				require.Len(t, accounts, 1)
				requireBodyMatchAccount(t, recorder.Body, accounts[0])
			},
		},
		{
			name: "BadCurrency",
			body: gin.H{
				"owner":    user.Username,
				"currency": "customCurrency",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeKey, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "Field validation for 'Currency' failed on the 'currency' tag")
			},
		},
		{
			name: "CreateSameAccCurrencyTwice",
			body: gin.H{
				"owner":    user.Username,
				"currency": util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeKey, user.Username, time.Minute)
				testStore.Store.CreateAccount(context.Background(), db.CreateAccountParams{
					Owner:    user.Username,
					Balance:  0,
					Currency: util.USD,
				})
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
			req, err := http.NewRequest(http.MethodPost, "/accounts", body)
			require.NoError(t, err)
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}

}

func TestGetAccountAPI(t *testing.T) {
	user, _ := db.CreateRandomUser(t, testServer.store)
	account := db.CreateRandomAccount(t, testServer.store, user)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeKey, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "ExpiredAuthToken",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeKey, user.Username, -time.Hour)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), "expired token")
			},
		},
		{
			name:      "NotFound-Unauthorized",
			accountID: account.ID + 3432423433,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeKey, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "ID less than 1",
			accountID: -account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeKey, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case-%d-%s", i, tc.name), func(t *testing.T) {
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	decoder := json.NewDecoder(body)
	var bodyAccount db.Account
	err := decoder.Decode(&bodyAccount)
	account.CreatedAt = account.CreatedAt.UTC()
	require.NoError(t, err)
	require.Equal(t, account, bodyAccount)
}
