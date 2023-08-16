package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/AdamDomagalsky/goes/2023/bank/util"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
)

var testDB *sql.DB
var testQueries *Queries

// go test main_test.go db.go
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load env config:", err)
	}

	testDB, err = sql.Open(config.DATABASE_DRVIER, util.DbURL(config))
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
