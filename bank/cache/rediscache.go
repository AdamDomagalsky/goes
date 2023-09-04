package kesh

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

var _ Cache = (*RedisCache)(nil)

type RedisCache struct {
	client *redis.Client
	cache  *cache.Cache
}

func NewRedisCache(redisOpt *redis.Options) (*RedisCache, error) {
	redisClient := redis.NewClient(redisOpt)

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	cacheClient := cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: nil,
	})

	return &RedisCache{
		client: redisClient,
		cache:  cacheClient,
	}, nil
}

func (rc *RedisCache) Set(rctx context.Context, key string, value interface{}, ttl time.Duration) (string, error) {
	return "", nil
}
func (rc *RedisCache) SetXX(rctx context.Context, key string, value interface{}, ttl time.Duration) bool {
	return true
}
func (rc *RedisCache) SetNX(rctx context.Context, key string, value interface{}, ttl time.Duration) bool {
	cmd := rc.client.SetNX(rctx, key, value, ttl)
	return cmd.Val()
}
func (rc *RedisCache) Get(rctx context.Context, key string) (string, error) {
	return rc.client.Get(rctx, key).Result()
}

// // or cache
//
//		// var obj string
//		// return obj, rc.cache.Get(rctx, key, &obj)
//	}
func (rc *RedisCache) Del(rctx context.Context, keys ...string) int64 { return -1 }

func (rc *RedisCache) Once(key string, value interface{}, ttl time.Duration) error {

	doWrapper := func(item *cache.Item) (interface{}, error) {

		return nil, nil
	}

	return rc.cache.Once(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   ttl,
		Do:    doWrapper,
	})
}
