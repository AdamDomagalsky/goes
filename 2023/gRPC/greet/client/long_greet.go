package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/greet/proto"
)

func doLongGreet(c pb.GreetServiceClient) {
	log.Printf("doGreet was invoked")

	reqs := []*pb.GreetRequest{
		{FirstName: "Foo"},
		{FirstName: "Baz"},
		{FirstName: "Bar"},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v\n", err)
	}
	for _, req := range reqs {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v\n", err)
	}

	log.Printf("LongGreet; %v\n", res.Result)
}
