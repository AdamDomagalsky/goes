package main

import (
	"context"
	"io"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
)

func doSum(c pb.CalculatorServiceClient) {
	log.Printf("doSum was invoked")
	res, err := c.Sum(context.Background(), &pb.SumRequest{
		A: 3,
		B: 10,
	})

	if err != nil {
		log.Fatalf("Could not sum: %v\n", err)
	}

	log.Printf("Sum(a,b): %v\n", res.Result)
}

func doPrimes(c pb.CalculatorServiceClient) {
	log.Printf("doPrimes was invoked")
	var N int64 = 120

	req := &pb.PrimeRequest{Number: N}
	stream, err := c.Primes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Primes: %v\n", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading the stream: %v\n", err)
		}
		log.Printf("Primes: %v\n", msg.Result)
	}
}
