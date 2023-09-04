package main

import (
	"context"
	"fmt"
	jejdis "goes/worker/cache"
	"time"

	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

func protoTest() {
	// https://gist.github.com/miguelmota/25568433ad8cfddb5ea556a5644c9fde
	userb4 := &pb.User{
		Username:          "test",
		Fullname:          "test",
		Email:             "test",
		PasswordChangedAt: nil,
		CreatedAt:         nil,
	}
	data, _ := proto.Marshal(userb4)
	// require.NoError(t, err)
	fmt.Println(data)
	userAfter := &pb.User{}
	_ = proto.Unmarshal(data, userAfter)
	// require.NoError(t, err)
	// require.True(t, proto.Equal(userb4, userAfter))
}

func getTest() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	for {

		err := rdb.Set(ctx, "key_kluczyk", "value_siemandero", 0).Err()
		if err != nil {
			panic(err)
		}

		val, err := rdb.Get(ctx, "key_kluczyk").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("key", val)

		time.Sleep(time.Millisecond * 2)

		val2, err := rdb.Get(ctx, "key2").Result()
		if err == redis.Nil {
			fmt.Println("key2 does not exist")
		} else if err != nil {
			panic(err)
		} else {
			fmt.Println("key2", val2)
		}

	}
}

func main() {
	jejdis.Kesz()

}
