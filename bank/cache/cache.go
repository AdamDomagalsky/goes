package kesh

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (string, error)
	SetXX(ctx context.Context, key string, value interface{}, ttl time.Duration) bool
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) bool
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) int64
	// Once(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}
