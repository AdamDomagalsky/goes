package main

import (
	"database/sql"
	"log"

	"github.com/AdamDomagalsky/goes/2023/bank/api"
	db "github.com/AdamDomagalsky/goes/2023/bank/db/sqlc"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://admin:secret@localhost:5432/pg-bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	var err error
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)

	if err != nil {
		log.Fatalf("Cannot start server: %+v\n ", err)
	}

}
