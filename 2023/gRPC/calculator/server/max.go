package main

import (
	"io"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
)

func (*Server) Max(stream pb.CalculatorService_MaxServer) error {
	log.Println("Max function was invoked")

	req, err := stream.Recv()
	if err == io.EOF {
		return nil

	}
	if err != nil {
		log.Fatalf("Error while reading client stream: %v\n", err)
	}

	max := req.Number
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}

		if req.Number > max {
			max = req.Number
		}
		stream.Send(&pb.MaxResponse{
			Result: max,
		})
	}
	return nil
}
