package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
)

func readBlog(c pb.BlogServiceClient, id string) *pb.Blog {
	log.Println("--invoked readBlog--")

	req := &pb.BlogId{Id: id}
	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Error happend while reading %v\n", err)
	}

	return res
}
