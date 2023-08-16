package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AdamDomagalsky/goes/2023/bank/api"
	db "github.com/AdamDomagalsky/goes/2023/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/2023/bank/util"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
)

func main() {
	config, err := util.LoadConfig(".")
	fmt.Printf("config: %v\n", config)
	if err != nil {
		log.Fatal("Cannot load env config:", err)
	}
	conn, err := sql.Open(config.DATABASE_DRVIER, util.DbURL(config))
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.SERVER_API_URL)
	if err != nil {
		log.Fatalf("Cannot start server: %+v\n ", err)
	}

}
