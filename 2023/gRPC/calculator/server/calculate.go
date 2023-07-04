package main

import (
	"context"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
)

func (s *Server) Sum(ctx context.Context, in *pb.CalculatorRequest) (*pb.CalculatorResponse, error) {

	log.Printf("Sum function was invoked with %v\n", in)

	return &pb.CalculatorResponse{
		Sum: int64(in.A) + int64(in.B),
	}, nil
}
