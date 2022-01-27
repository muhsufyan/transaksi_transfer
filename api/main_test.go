package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/muhsufyan/transaksi_transfer/db/sqlc"
	"github.com/muhsufyan/transaksi_transfer/util"
	"github.com/stretchr/testify/require"
)

// tambhn ada token
func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	// change gin test mode
	gin.SetMode(gin.TestMode)
	//running unit test
	os.Exit(m.Run())
}
