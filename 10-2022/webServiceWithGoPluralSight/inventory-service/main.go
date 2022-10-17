package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // side effect for db, not calling implicitly
	"goes/webServiceWithGoPluralSight/database"
	"goes/webServiceWithGoPluralSight/product"
	"goes/webServiceWithGoPluralSight/receipt"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoute(apiBasePath)
	receipt.SetupRoute(apiBasePath)
	log.Fatalln(http.ListenAndServe(":5005", nil)) // Go uses the DefaultServeMux to handle incoming requests.
}
