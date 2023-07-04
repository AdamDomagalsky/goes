package main

import (
	"context"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/greet/proto"
)

func (s *Server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {

	log.Printf("Greeting function was invoked with %v\n", in)

	return &pb.GreetResponse{
		Result: "Hello" + in.FirstName,
	}, nil
}
