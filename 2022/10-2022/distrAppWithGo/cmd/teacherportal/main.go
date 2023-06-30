package main

import (
	"context"
	"fmt"
	stdlog "log"

	"goes/distrAppWithGo/log"
	"goes/distrAppWithGo/registry"
	"goes/distrAppWithGo/service"
	"goes/distrAppWithGo/teacherportal"
)

func main() {
	err := teacherportal.ImportTemplates()
	if err != nil {
		stdlog.Fatal(err)
	}

	host, portal := "localhost", "9000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, portal)
	var r registry.Registration
	r.ServiceName = registry.TeacherPortal
	r.ServiceURL = serviceAddress
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
		registry.GradingService,
	}
	r.ServiceUpdateURL = r.ServiceURL + "/services"
	r.HeartbeatURL = r.ServiceURL + "/heartbeat"
	ctx, err := service.Start(context.Background(),
		host,
		portal,
		r,
		teacherportal.RegisterHandlers)
	if err != nil {
		stdlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("logging servcies found at: %v\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	stdlog.Printf("test if stdlog is beeing send to app.log")
	<-ctx.Done()
	fmt.Println("Shutting down teacher portal")
}
