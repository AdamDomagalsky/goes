package main

import (
	"context"
	"log"
	"net"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50051"
var uri string = "mongodb://root:root@localhost:27017/"

var collection *mongo.Collection

type Server struct {
	pb.BlogServiceServer
}

func main() {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("blogdb").Collection("blog")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}
	log.Printf("Listening on %s\n", addr)

	s := grpc.NewServer()
	pb.RegisterBlogServiceServer(s, &Server{})
	// reflection.Register(s) // required to make evans work
	if err != s.Serve(lis) {
		log.Fatalf("Failed to serve %v\n", err)
	}
}
