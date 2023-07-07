package main

import (
	"context"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func doSqrt(c pb.CalculatorServiceClient, number float64) {

	res, err := c.Sqrt(context.Background(), &pb.SqrtRequest{Number: number})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {

			switch e.Code() {
			case codes.InvalidArgument:
				log.Printf("%v", e.Message())
			default:
				log.Printf("Code(%v) Message: %s\n", e.Code(), e.Message())
			}
			return
		} else {
			log.Fatalf("A non gRPC error: %v\n", err)
		}
	}

	log.Printf("Sqrt(%v): %+v\n", number, res.Result)
}
