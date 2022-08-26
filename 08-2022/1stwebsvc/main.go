package main

import (
	controllers "first-websvc/controlles"
	"fmt"
	"net/http"
)

func main() {

	controllers.RegisterControllers()
	http.ListenAndServe(":3000", nil)

	fmt.Println("siema")
}
