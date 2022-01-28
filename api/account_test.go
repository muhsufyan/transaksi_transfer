package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/muhsufyan/transaksi_transfer/db/mock"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
	"github.com/muhsufyan/transaksi_transfer/token"
	"github.com/muhsufyan/transaksi_transfer/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	// create new random user
	user, _ := randomUser(t)
	// create new account with random generate
	account := randomAccount(user.Username)
	// 1st ttd : LIST OF TEST CASE. anonimous class u/ menyimpan test data
	testCases := []struct {
		// setiap test case have unik name
		name string
		// id akun yg ingin didptkan
		accountID int64
		// setup header otorisasi dr request. param 3 token maker interface to buat token akses
		setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		// GetAccount stubs u/ each skenario will be build differently. MockStore to build the stub karena suite untuk tujuan dr setiap test case
		buildStubs func(store *mockdb.MockStore)
		// cek output dr API
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			// skenario happy test (test data from response body)
			name:      "OK",
			accountID: account.ID,
			// implement otorisasi
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// setup auth func. waktunya 1 menit
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs untuk this mock store
				// GetAccount is interface  & ada di db/querier.go
				// run 1 kali
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// cek response
				require.Equal(t, http.StatusOK, recorder.Code)
				// cek response body
				// response body tersimpan in recorder.Body(param 2), generated account (param 3)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// test case token usernya tdk sama dg dirinya
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			// implement otorisasi
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// setup auth func.
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs untuk this mock store
				// GetAccount is interface  & ada di db/querier.go
				// run 1 kali
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// cek response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// test case token client tdk punya token akses
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			// implement otorisasi
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs untuk this mock store
				// expectednya GetAccount tdk dijlnkan
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// cek response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// test when account is not found, expected Not Found
		{
			name:      "NotFound",
			accountID: account.ID,
			// implement otorisasi
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// setup auth func. waktunya 1 menit
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs untuk this mock store
				// GetAccount is interface  & ada di db/querier.go
				// run 1 kali
				// returned 1 akun kosong & karena akun tdk ada maka not found dg sql.ErrNoRows
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// cek response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		// test ttd untuk internal server error, expected internal error
		{
			name:      "InternalError",
			accountID: account.ID,
			// implement otorisasi
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// setup auth func. waktunya 1 menit
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// returned 1 akun kosong & karena akun tdk ada maka not found dg sql.ErrNoRows
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// cek response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		// test ttd untuk bad request (user kirim param yg invalid ke API), expected bad request
		// skenario nya tdk terpenuhi binding dimana id nya 0 menyebabkan invalid id shrsnya >= 1
		{
			name:      "InvalidID",
			accountID: 0,
			// implement otorisasi
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// setup auth func. waktunya 1 menit
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// GetAccount param 2 since id invalid GetAccount shrsnya tdk dipanggil oleh handler
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// cek response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// TODO add more test cases
	}
	// jlnkan semua test case dg loop
	for i := range testCases {
		// simpan data dr current test case
		tc := testCases[i]
		// run each case as a separate sub-test of this unit test. tc.name : nama test case
		t.Run(tc.name, func(t *testing.T) {
			// gomock controller (ada di mock/store.go) as input/param
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			// start test server & send request
			server := newTestServer(t, store)
			// we not use real http api tp use record feature dr httptest
			recorder := httptest.NewRecorder()
			// api yg ingin kita panggil. data ID untuk setiap test case berbeda"
			url := fmt.Sprintf("/account/%d", tc.accountID)
			// request body nya nil
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			// create obj recorder & request. ini akan send request melalui server router & response berupa record berasal dr recorder
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// generate random akun
func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

// CEK RESPONSE BODY
// param 2 : response body, param 3 : obj akun untuk dibandingkan
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// read all data from response body (data dr param response body)
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	// to menyimpan obj akun got from response body data
	var gotAccount db.Account
	// unmarshal data to the gotAccount obj
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func TestCreateAccountAPI(t *testing.T) {
	// create new random user
	user, _ := randomUser(t)
	// create new account with random generate
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"owner":    account.Owner,
				"currency": "invalid",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	// create new random user
	user, _ := randomUser(t)
	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					// Owner:  user.Username,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/account"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
