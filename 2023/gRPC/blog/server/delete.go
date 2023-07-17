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

func (*Server) DeleteBlog(ctx context.Context, in *pb.BlogId) (*emptypb.Empty, error) {
	fmt.Printf("DeleteBlog was invoked with %v\n", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot parse Id")
	}
	res, err := collection.DeleteOne(
		ctx,
		bson.M{"_id": oid},
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot delete object in MongoDB: %v\n", err)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Blog not found")
	}

	return &emptypb.Empty{}, nil
}
