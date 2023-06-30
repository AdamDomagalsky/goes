// Package main is typically used to define a single, executable command and its associated logic.
//
// Functions in the main package should normally avoid containing business logic; they
// should focus on initializing the application and handling any errors that are returned by
// any library calls.
package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/pluralsight/libmanager/services/events"
)

func init() {
	fmt.Println("init in main.go")
}

// The main function, when defined in a main package, defines
// the entry and exit point for an executable command.
func main() {
	log.Println("Starting all services")

	// start the services
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go startServer(wg, events.StartServer)
	wg.Wait()
	log.Println("All services stopped")
}

// startServer is a function in the main package that supports the command.
func startServer(wg *sync.WaitGroup, startFunc func() error) {
	err := startFunc()
	wg.Done()
	if err != nil {
		log.Fatal(err)
	}
}
