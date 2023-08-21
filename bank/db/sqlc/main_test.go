package db

import (
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
)

var (
	testEnv *TestStore
)

// go test main_test.go db.go
func TestMain(m *testing.M) {

	testEnv = NewStoreTestStore("../..", testEnv)
	exitCode := m.Run()
	testEnv.Teardown()

	os.Exit(exitCode)
}
