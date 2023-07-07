package main

import (
	"io"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/greet/proto"
)

func (*Server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	log.Println("LongGreet function was invoked")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}
		res := "Helo " + req.FirstName + "!"
		err = stream.Send(&pb.GreetResponse{
			Result: res,
		})
		if err != nil {
			log.Fatalf("Error while sending responde to the stream: %v\n", err)
		}
	}

	return nil
}
