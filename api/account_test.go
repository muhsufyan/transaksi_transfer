package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"bytes"
	"testing"
	"net/http"
	"net/http/httptest"
	"database/sql"
	mockdb "github.com/muhsufyan/transaksi_transfer/db/mock"
	"github.com/golang/mock/gomock"
	db"github.com/muhsufyan/transaksi_transfer/db/sqlc"
	"github.com/muhsufyan/transaksi_transfer/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	// create new account with random generate
	account := randomAccount()
	// 1st ttd : LIST OF TEST CASE. anonimous class u/ menyimpan test data
	testCases :=  []struct{
		// setiap test case have unik name 
		name string
		// id akun yg ingin didptkan
		accountID int64
		// GetAccount stubs u/ each skenario will be build differently. MockStore to build the stub karena suite untuk tujuan dr setiap test case
		buildStubs func(store *mockdb.MockStore)
		// cek output dr API
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
		
	}{
		{
			// skenario happy test (test data from response body)
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore){
				// build stubs untuk this mock store
				// GetAccount is interface  & ada di db/querier.go
				// run 1 kali
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				// cek response
				require.Equal(t, http.StatusOK, recorder.Code)
				// cek response body
				// response body tersimpan in recorder.Body(param 2), generated account (param 3) 
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// test when account is not found, expected Not Found
		{
			name: "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore){
				// build stubs untuk this mock store
				// GetAccount is interface  & ada di db/querier.go
				// run 1 kali
				// returned 1 akun kosong & karena akun tdk ada maka not found dg sql.ErrNoRows
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				// cek response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		// TODO add more test cases
	}
	// jlnkan semua test case dg loop
	for i := range testCases{
		// simpan data dr current test case
		tc := testCases[i]
		// run each case as a separate sub-test of this unit test. tc.name : nama test case
		t.Run(tc.name, func(t *testing.T){
			// gomock controller (ada di mock/store.go) as input/param
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			// start test server & send request
			server := NewServer(store)
			// we not use real http api tp use record feature dr httptest
			recorder := httptest.NewRecorder()
			// api yg ingin kita panggil
			url := fmt.Sprintf("/account/%d", account.ID)
			// request body nya nil
			request, err := http.NewRequest(http.MethodGet, url, nil) 
			require.NoError(t, err)
			// create obj recorder & request. ini akan send request melalui server router & response berupa record berasal dr recorder
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
}
}

// generate random akun
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
// CEK RESPONSE BODY
// param 2 : response body, param 3 : obj akun untuk dibandingkan
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account)  {
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