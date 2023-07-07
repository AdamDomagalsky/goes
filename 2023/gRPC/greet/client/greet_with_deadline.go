package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/greet/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func doGreetWithDeadline(c pb.GreetServiceClient, timeout time.Duration) {
	log.Printf("doGreetWithDeadline was invoked with client timeout %v\n", timeout)
	// c.GreetWithDeadline(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &pb.GreetRequest{
		FirstName: "Deadlinerzz",
	}
	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			if e.Code() == codes.DeadlineExceeded {
				log.Printf("Deadline exceeded %v\n", timeout)
				return
			} else {
				log.Fatalf("gRPC error: %v", err)

			}
		} else {
			log.Fatalf("Not gRPC error: %v", err)
		}
	}

	fmt.Printf("GreetWithDeadline: %s\n", res.Result)
}
