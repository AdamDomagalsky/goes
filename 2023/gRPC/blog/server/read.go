package main

import (
	"context"
	"fmt"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*Server) ReadBlog(ctx context.Context, in *pb.BlogId) (*pb.Blog, error) {
	fmt.Printf("ReadBlog was invoked with %v\n", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	data := &BlogItem{}
	res := collection.FindOne(ctx, bson.M{"_id": oid})
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot find blog with the ID provided",
		)
	}

	return documentToBlog(data), nil
}

// func (*Server) UpdateBlog(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error)
// func (*Server) DeleteBlog(ctx context.Context, in *pb.BlogId) (*emptypb.Empty, error)
// func (*Server) ListBlogs(ctx context.Context, in *emptypb.Empty) (pb.BlogService_ListBlogsClient, error)
