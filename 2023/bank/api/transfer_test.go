package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/AdamDomagalsky/goes/2023/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/2023/bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	amount := int64(10)
	server := testServer
	store := server.store
	user1, _ := db.CreateRandomUser(t, store)
	user2, _ := db.CreateRandomUser(t, store)
	user3, _ := db.CreateRandomUser(t, store)
	account1, err := store.CreateAccount(context.Background(), db.CreateAccountParams{
		Owner:    user1.Username,
		Currency: util.USD,
		Balance:  1000,
	})
	require.NoError(t, err)

	account2, _ := store.CreateAccount(context.Background(), db.CreateAccountParams{
		Owner:    user2.Username,
		Currency: util.USD,
		Balance:  0,
	})
	require.NoError(t, err)

	account3, _ := store.CreateAccount(context.Background(), db.CreateAccountParams{
		Owner:    user3.Username,
		Currency: util.EUR,
		Balance:  1000,
	})
	require.NoError(t, err)

	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account_id": 100001,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   100001,
				"amount":          amount,
				"currency":        util.USD,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "BadRequestAmountAccountsIDsLessThanOne",
			body: gin.H{
				"from_account_id": -account1.ID,
				"to_account_id":   -account2.ID,
				"amount":          -amount,
				"currency":        util.USD,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "Error:Field validation for 'Amount' failed on the 'gt' tag")
				require.Contains(t, recorder.Body.String(), "Error:Field validation for 'FromAccountID' failed on the 'min' tag")
				require.Contains(t, recorder.Body.String(), "Error:Field validation for 'ToAccountID' failed on the 'min' tag")
			},
		},
		{
			name: "CurrencyMismatch",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        util.EUR,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "currency mismatch EUR vs USD")
			},
		},
		{
			name: "todo",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        util.EUR,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "currency mismatch EUR vs USD")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(fmt.Sprintf("case-%d-%s", i, tc.name), func(t *testing.T) {
			recorder := httptest.NewRecorder()
			url := "/transfers"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}

}
