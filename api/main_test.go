package api

import (
	"os"
	"testing"
	"github.com/gin-gonic/gin"
)
func TestMain(m *testing.M) {
	// change gin test mode
	gin.SetMode(gin.TestMode)
	//running unit test
	os.Exit(m.Run())
}
