package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/greet/proto"
)

func doGreetEveryone(c pb.GreetServiceClient) {

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	names := []string{
		"Adam",
		"Pablo",
		"Krzysiu",
	}

	go func(s pb.GreetService_GreetEveryoneClient, nms []string) {
		for _, name := range nms {
			log.Printf("Sending requst: %v\n", name)
			s.Send(&pb.GreetRequest{
				FirstName: name,
			})

			time.Sleep(1 * time.Second)
		}
		s.CloseSend()
	}(stream, names)

	waitc := make(chan struct{})
	go func(wait_channel chan struct{}, s pb.GreetService_GreetEveryoneClient) {
		for {
			res, err := s.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error while receiving: %v\n", err)
				break
			}

			fmt.Printf("Received: %v\n", res.Result)
		}
		close(wait_channel)
	}(waitc, stream)
	<-waitc
}
