package api

import (
	"os"
	"testing"

	db "github.com/AdamDomagalsky/goes/2023/bank/db/sqlc"
	"github.com/gin-gonic/gin"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
)

var (
	testStore  *db.TestStore
	testServer *Server
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testStore = db.NewStoreTestStore("..", testStore)
	testServer = NewServer(testStore.Store)
	exitCode := m.Run()
	testStore.Teardown()
	os.Exit(exitCode)
}
