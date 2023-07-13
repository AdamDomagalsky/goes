package main

import (
	"log"
	"net"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var addr string = "0.0.0.0:50051"

type Server struct {
	pb.CalculatorServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}
	log.Printf("Listening on %s\n", addr)

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &Server{})
	reflection.Register(s) // required to make evans work
	if err != s.Serve(lis) {
		log.Fatalf("Failed to serve %v\n", err)
	}
}
