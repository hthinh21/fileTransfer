package store

import (
	"context"
	"time"

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

func (r *RedisStore) Save(code string, path string) error {
	return r.client.Set(r.ctx, code, path, 10*time.Minute).Err()
}

func (r *RedisStore) Get(code string) (string, error) {
	return r.client.Get(r.ctx, code).Result()
}