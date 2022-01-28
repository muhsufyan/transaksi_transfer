package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhsufyan/transaksi_transfer/token"
	"github.com/stretchr/testify/require"
)

// func ini will used in banyak kasus dg config yg berbeda"
func addAuthorization(t *testing.T, request *http.Request, tokenMaker token.Maker, authorizationType string, username string, duration time.Duration) {
	// buat token baru
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	// buat header otorisasi, param 2 akses token
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	// set header dr request
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
func TestAuthMiddleware(t *testing.T) {
	/*
		memakai ttd/table-driven test
	*/
	// all test case store in anonymous struct
	testCases := []struct {
		// setiap test diberi nama
		name string
		// setup header otorisasi dr request. param 3 token maker interface to buat token akses
		setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		// setiap test hrs punya func untuk cek response
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// disini all test case dibuat
		// case ok/sukses
		{
			name: "OK",
			// implement of interface
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// buat new token akses
				// add akses token ke header otorisasi dari request

				// setup auth func. username : user dan waktunya 1 menit
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			// setiap test hrs punya func untuk cek response
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		// case tdk ada otorisasi (client tdk punya header otorisasi)
		{
			name: "NoAuthorization",
			// implement of interface
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			// setiap test hrs punya func untuk cek response
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// case otorisasi tdk didukung server. server dpt memiliki banyak jenis otorisasi sprti oauth, token, keyapi dll)
		{
			name: "UnsupportedAuthorization",
			// implement of interface
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// setup auth func.jenis otorisasinya error karena tdk didukung
				addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)
			},
			// setiap test hrs punya func untuk cek response
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// case client tdk ada tipe prefix pd header otorisasinya
		{
			name: "InvalidAuthorizationFormat",
			// implement of interface
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// setup auth func.jenis otorisasinya error(kosongkan) karena tdk didukung
				addAuthorization(t, request, tokenMaker, "", "user", time.Minute)
			},
			// setiap test hrs punya func untuk cek response
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// case token client kadaluarsa
		{
			name: "TokenExpired",
			// implement of interface
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// setup auth func.
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", -time.Minute)
			},
			// setiap test hrs punya func untuk cek response
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	// iterate test case via testCases slice
	for i := range testCases {
		// store current test case
		tc := testCases[i]
		// generate sub-test
		t.Run(tc.name, func(t *testing.T) {
			// callback func ini berisi content utama dr sub-test

			// buat new test server. in middleware test we tdk perlu mengakses Store so param 2(mrpkn db.Store) nil
			server := newTestServer(t, nil)
			// url. add simple api route & handler u/ testing middleware
			authPath := "/auth"
			// jlnkan server
			server.router.GET(
				// url
				authPath,
				authMiddleware(server.tokenMaker),
				// handler func
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)
			// send request to api. to record the call
			recorder := httptest.NewRecorder()
			// buat new request. param 3(request body)
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			//to add header otorisasi to the request
			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			// verify result
			tc.checkResponse(t, recorder)
		})
	}
}
