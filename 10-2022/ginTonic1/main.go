package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello world!")
	})

	log.Fatal(router.Run(":3000"))
}
