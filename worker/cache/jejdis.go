package jejdis

import (
	"fmt"
	"time"

	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
}

func NewRedis() *Redis {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &Redis{
		rdb: rdb,
	}
}

func Kesz() {

	r := NewRedis()
	kesz := cache.New(&cache.Options{
		Redis: r.rdb,
	})

	obj := &pb.User{}

	doIfNotInCache := func(*cache.Item) (interface{}, error) {
		userb4 := &pb.User{
			Username:          "DO",
			Fullname:          "DO",
			Email:             "DO",
			PasswordChangedAt: nil,
			CreatedAt:         nil,
		}
		fmt.Println(userb4)

		fmt.Println("gej3")

		return userb4, nil
	}

	err := kesz.Once(&cache.Item{
		Key:   "gej4",
		Value: obj,
		TTL:   time.Second * 10,
		Do:    doIfNotInCache,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(obj)
}
