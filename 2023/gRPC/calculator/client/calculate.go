package main

import (
	"context"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
)

func doSum(c pb.CalculatorServiceClient) {
	log.Printf("doSum was invoked")
	res, err := c.Sum(context.Background(), &pb.CalculatorRequest{
		A: 3,
		B: 10,
	})

	if err != nil {
		log.Fatalf("Could not sum: %v\n", err)
	}

	log.Printf("Sum(a,b): %v\n", res.Sum)
}
