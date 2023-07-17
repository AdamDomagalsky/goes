package main

import (
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "0.0.0.0:50051"

func main() {
	tls := false
	options := []grpc.DialOption{}
	if tls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalf("Error while loading CA trust certifcate: %v\n", err)
		}
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.Dial(addr, options...)
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewBlogServiceClient(conn)
	id := createBlog(c)
	log.Printf("Blog has been created id: %v\n", id)
	log.Println(readBlog(c, id))
	// log.Println(readBlog(c, "anot2137"))
	updateBlog(c, id)
	listBlog(c)
	deleteBlog(c, id)
	listBlog(c)
}
