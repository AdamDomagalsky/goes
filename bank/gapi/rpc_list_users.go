package gapi

import (
	"context"
	"fmt"
	"time"

	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) ListUsers(ctx context.Context, _ *emptypb.Empty) (*pb.ListUsersResponse, error) {
	// echo "GET http://localhost:2138/v1/list_users" | vegeta attack -duration=5s | tee results2.bin | vegeta report
	listUserProtoResp := &pb.ListUsersResponse{}
	const redisKey = "users:ALL"

	cachedUsersList, err := server.cache.Get(ctx, redisKey)
	if err != nil {
		us, _ := server.store.ListUsers(ctx)
		for _, u := range us {
			listUserProtoResp.User = append(listUserProtoResp.User, convertUser(u))
		}
		bajts, _ := proto.Marshal(listUserProtoResp)
		go server.cache.SetNX(ctx, redisKey, bajts, time.Second*1)
		fmt.Print("cache MISS all_users, getting from db and setting in cache")

	} else {
		proto.Unmarshal([]byte(cachedUsersList), listUserProtoResp)
		fmt.Print("cache HIT all_users, getting from cache!")
	}

	return listUserProtoResp, nil
}
