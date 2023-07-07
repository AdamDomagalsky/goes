package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
)

func doMax(c pb.CalculatorServiceClient) {
	stream, err := c.Max(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream%v\n", err)
	}
	reqs := []int{-300, -500, 0, 5, 321, -456, 666, -999}
	go func(s pb.CalculatorService_MaxClient, numbers []int) {
		for _, number := range numbers {
			log.Printf("Sending requst: %v\n", number)
			s.Send(&pb.MaxRequest{Number: int64(number)})
			time.Sleep(1 * time.Second)

		}
		s.CloseSend()
	}(stream, reqs)

	wait_channel := make(chan interface{})
	go func(waitc chan interface{}, s pb.CalculatorService_MaxClient) {

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

		close(waitc)
	}(wait_channel, stream)
	<-wait_channel
}
