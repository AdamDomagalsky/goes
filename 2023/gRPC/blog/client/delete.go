package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
)

func deleteBlog(c pb.BlogServiceClient, id string) {
	log.Println("--invoked deleteBlog--")

	req := &pb.BlogId{
		Id: id,
	}
	_, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Error happend while deleteing %v\n", err)
	}
	fmt.Printf("blog deleted, id: %v\n", req)
}
