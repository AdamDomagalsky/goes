package main

import (
	"context"
	"log"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
)

func createBlog(c pb.BlogServiceClient) string {
	log.Println("--invoked createBlog--")
	data, err := c.CreateBlog(context.Background(), &pb.Blog{
		AuthorId: "Stefano Terantzino",
		Title:    "Dzis w klubie bedzie beng",
		Content:  "Kontent eyy",
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	return data.Id
}
