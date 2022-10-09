package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"goes/distrAppWithGo/registry"
)

func main() {
	http.Handle("/services", &registry.RegistryService{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()
	go func() {
		fmt.Println("Registry service started. Press any key to stop.")
		var w8 string
		_, _ = fmt.Scanln(&w8)
		_ = srv.Shutdown(ctx)
		cancel()
	}()
	<-ctx.Done()
	fmt.Println("Shutting down registry service")
}
