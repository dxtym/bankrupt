package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/dxtym/bankrupt/db/mock"
	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTransferAPI(t *testing.T) {
	amount := int64(10)

	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()
	
	account1.Currency = utils.USD
	account2.Currency = utils.USD
	account3.Currency = utils.EUR

	testCases := []struct{
		name string
		body gin.H
		buildStubs func(s *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": utils.USD,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).Return(account1, nil)
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountId: account1.ID,
					ToAccountId: account2.ID,
					Amount: amount,
				}
				s.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).AnyTimes()
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": utils.USD,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).Return(db.Account{}, sql.ErrNoRows)
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				s.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": utils.USD,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).Return(account1, nil)
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).Return(db.Account{}, sql.ErrNoRows)
				s.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": account3.ID,
				"to_account_id": account1.ID,
				"amount": amount,
				"currency": utils.USD,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
					Times(1).Return(account3, nil)
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(0)
				s.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)	
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account3.ID,
				"amount": amount,
				"currency": utils.USD,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).Return(account1, nil)
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
					Times(1).Return(account3, nil)
				s.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)	
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": "XYZ",
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).Times(0)
				s.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)	
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": utils.USD,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Account{}, sql.ErrConnDone)
				s.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)	
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": -amount,
				"currency": utils.USD,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).Times(0)
				s.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)	
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TransferTxError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": utils.USD,
			},
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).Return(account1, nil)
				s.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).Return(account2, nil)
				s.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}