package store

import (
	"context"
	"encoding/json"
	"time"

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

func NewRedisStore() *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}
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
