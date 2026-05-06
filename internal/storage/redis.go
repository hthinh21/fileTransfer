package storage

import (
	"context"
	"encoding/json"
	"time"
	"os"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type FileMeta struct {
	ObjectKey string `json:"object_key"`
	FileName  string `json:"file_name"`
}

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore() (*RedisStore, error) {
    addr := os.Getenv("REDIS_ADDR")
    if addr == "" {
        addr = "localhost:6379"
    }
    client := redis.NewClient(&redis.Options{Addr: addr})

    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, fmt.Errorf("redis connect failed: %w", err)
    }

    return &RedisStore{
        client: client,
        ctx:    context.Background(),
    }, nil
}

func (r *RedisStore) Save(code string, meta FileMeta) error {
	payload, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, code, payload, 10*time.Minute).Err()
}

func (r *RedisStore) Get(code string) (FileMeta, error) {
	payload, err := r.client.Get(r.ctx, code).Result()
	if err != nil {
		return FileMeta{}, err
	}

	var meta FileMeta
	if err := json.Unmarshal([]byte(payload), &meta); err != nil {
		return FileMeta{}, err
	}

	return meta, nil
}

func (r *RedisStore) Delete(code string) error {
	return r.client.Del(r.ctx, code).Err()
}
