package main

import (
	"context"
	"io"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func listBlog(c pb.BlogServiceClient) {
	log.Println("--invoked listBlog--")

	stream, err := c.ListBlogs(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Smthing happened %v\n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Smthing happened %v\n", err)
		}

		log.Println(res)
	}
}
