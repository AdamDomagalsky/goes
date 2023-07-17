package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
)

func updateBlog(c pb.BlogServiceClient, id string) {
	log.Println("--invoked updateBlog--")

	req := &pb.Blog{
		Id:       id,
		Title:    "_",
		Content:  "_",
		AuthorId: "New Author",
	}
	_, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Error happend while updating %v\n", err)
	}
}
