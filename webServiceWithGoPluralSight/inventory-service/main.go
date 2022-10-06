package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // side effect for db, not calling implicitly
	"goes/webServiceWithGoPluralSight/database"
	"goes/webServiceWithGoPluralSight/product"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoute(apiBasePath)
	log.Fatalln(http.ListenAndServe(":5005", nil))
}
