package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/AdamDomagalsky/goes/2023/bank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := db.CreateRandomAccount(t, testServer.store)

	testCases := []struct {
		name      string
		accountID int64
		checks    func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			checks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID + 3432423433,
			checks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "ID less than 1",
			accountID: -account.ID,
			checks: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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
			testServer.router.ServeHTTP(recorder, req)

			tc.checks(t, recorder)
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
