package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/AdamDomagalsky/goes/2023/bank/util"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://admin:secret@localhost:5432/pg-bank?sslmode=disable"
)

// go test main_test.go db.go
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load env config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DATABASE_URL)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
