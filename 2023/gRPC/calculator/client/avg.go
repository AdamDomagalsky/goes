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

	reqs := []*pb.AvgRequest{
		{Number: 10},
		{Number: 10},
		{Number: 10},
		{Number: 10},
	}

	for _, req := range reqs {
		log.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from Avg: %v\n", err)
	}

	log.Printf("Avg; %v\n", res.Result)

}
