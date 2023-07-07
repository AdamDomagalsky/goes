package main

import (
	"context"
	"fmt"
	"log"
	"math"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*Server) Sqrt(c context.Context, in *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	log.Printf("Sqrt function was invoked with %v\n", in)

	number := in.Number
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Negative Real Square root of %v does not exist", number))
	}

	return &pb.SqrtResponse{
		Result: math.Sqrt(number),
	}, nil
}
