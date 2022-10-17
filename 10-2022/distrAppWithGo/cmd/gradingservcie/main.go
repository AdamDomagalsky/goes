package main

import (
	"context"
	"fmt"
	stdlog "log"

	"goes/distrAppWithGo/grades"
	"goes/distrAppWithGo/log"
	"goes/distrAppWithGo/registry"
	"goes/distrAppWithGo/service"
)

func main() {
	log.Run("./app.log")
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
	var r registry.Registration
	r.ServiceName = registry.GradingService
	r.ServiceURL = serviceAddress
	r.RequiredServices = []registry.ServiceName{registry.LogService}
	r.ServiceUpdateURL = r.ServiceURL + "/services"
	r.HeartbeatURL = r.ServiceURL + "/heartbeat"
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers)
	if err != nil {
		stdlog.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("logging servcies found at: %v\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	stdlog.Printf("test if stdlog is beeing send to app.log")

	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
