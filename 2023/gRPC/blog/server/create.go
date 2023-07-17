package main

import (
	"context"
	"fmt"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*Server) CreateBlog(ctx context.Context, in *pb.Blog) (*pb.BlogId, error) {
	fmt.Printf("CreateBlog was invoked with %v\n", in)
	res, err := collection.InsertOne(ctx, &BlogItem{
		AuthorId: in.AuthorId,
		Content:  in.Content,
		Title:    in.Title,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal Error: %v\n", err), err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "Cannot convert to OID")
	}

	return &pb.BlogId{
		Id: oid.Hex(),
	}, nil
}
