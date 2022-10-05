package main

import (
	"log"
	"net/http"

	"goes/webServiceWithGoPluralSight/product"
)

const apiBasePath = "/api"

func main() {
	product.SetupRoute(apiBasePath)
	log.Fatalln(http.ListenAndServe(":5005", nil))
}
