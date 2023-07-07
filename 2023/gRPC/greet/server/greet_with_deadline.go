package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/greet/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*Server) GreetWithDeadline(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Printf("GreetWithDeadline with %v\n", in)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			msg := "The client cancel the request!"
			fmt.Println(msg)
			return nil, status.Error(codes.Canceled, msg)
		}
		time.Sleep(1 * time.Second)
	}

	return &pb.GreetResponse{
		Result: "Hello " + in.FirstName + "!",
	}, nil
}
