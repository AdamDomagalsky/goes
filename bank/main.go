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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	if config.GIN_MODE != "release" {
		err = db.MigrateUp(conn, config.DATABASE_NAME)
		if err != nil {
			if err.Error() != "no change" {
				log.Fatal(
					fmt.Sprintf("GIN-%s, MigratingUp(%s) - failed:", config.GIN_MODE, config.DATABASE_NAME),
					err)
			} else {
				log.Printf("GIN-%s, MigratingUp(%s) - no change\n", config.GIN_MODE, config.DATABASE_NAME)
			}
		} else {
			log.Printf("GIN-%s, MigratingUp(%s) - succeed \n", config.GIN_MODE, config.DATABASE_NAME)
		}
	}

	err = server.Start(config.SERVER_API_URL)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
