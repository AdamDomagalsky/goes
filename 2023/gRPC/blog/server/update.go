package main

import (
	"context"
	"fmt"

	pb "github.com/AdamDomagalsky/goes/2023/gRPC/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (*Server) UpdateBlog(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error) {
	fmt.Printf("UpdateBlog was invoked with %v\n", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": oid}, // filter
		bson.M{"$set": bson.M{"author_id": in.AuthorId}}, // change only 1 filed author_id
	)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Could not update",
		)
	}

	if res.MatchedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot find blog with Id",
		)
	}

	return &emptypb.Empty{}, nil

}

// func (*Server) UpdateBlog(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error)
// func (*Server) DeleteBlog(ctx context.Context, in *pb.BlogId) (*emptypb.Empty, error)
// func (*Server) ListBlogs(ctx context.Context, in *emptypb.Empty) (pb.BlogService_ListBlogsClient, error)
