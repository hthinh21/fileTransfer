package store

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore() *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}
}