package kesh

import (
	"context"
	"time"
)

var _ Cache = (*NullCache)(nil)

type NullCache struct{}

func NewNullCache() (*NullCache, error) {
	return &NullCache{}, nil
}

func (c *NullCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (string, error) {
	return "", nil
}
func (c *NullCache) SetXX(ctx context.Context, key string, value interface{}, ttl time.Duration) bool {
	return true
}
func (c *NullCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) bool {
	return true
}
func (c *NullCache) Get(ctx context.Context, key string) (string, error)         { return "", nil }
func (c *NullCache) Del(ctx context.Context, keys ...string) int64               { return -1 }
func (c *NullCache) Once(key string, value interface{}, ttl time.Duration) error { return nil }
