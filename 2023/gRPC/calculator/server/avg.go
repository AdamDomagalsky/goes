package main

import (
	"io"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
)

func (*Server) Avg(stream pb.CalculatorService_AvgServer) error {

	sum := 0
	for i := 0; ; i++ {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AvgResponse{
				Result: (float64(sum) / float64(i)),
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}

		log.Printf("Receiving: %v\n", req)
		sum += int(req.Number)
	}
}
