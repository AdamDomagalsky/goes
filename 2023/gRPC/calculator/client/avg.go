package main

import (
	"context"
	"log"
	"time"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
)

func doAvg(c pb.CalculatorServiceClient) {
	log.Printf("doAvg was invoked")

	stream, err := c.Avg(context.Background())
	if err != nil {
		log.Fatalf("Error while calling doAvg: %v\n", err)
	}

	numbers := []int32{1, 2, 3, 4, 5, 6, 7, 8}
	for _, number := range numbers {
		log.Printf("Sending req: %v\n", number)
		stream.Send(&pb.AvgRequest{Number: uint64(number)})
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from Avg: %v\n", err)
	}

	log.Printf("Avg; %v\n", res.Result)

}
