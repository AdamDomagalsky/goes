package main

import (
	"log"
	"net"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr string = "0.0.0.0:50051"

type Server struct {
	pb.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}
	log.Printf("Listening on %s\n", addr)
	tls := true
	options := []grpc.ServerOption{}
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Fail loading certifcates: %v\n", err)
		}
		options = append(options, grpc.Creds(creds))
	}
	s := grpc.NewServer(options...)
	pb.RegisterGreetServiceServer(s, &Server{})
	if err != s.Serve(lis) {
		log.Fatalf("Failed to serve %v\n", err)
	}
}
