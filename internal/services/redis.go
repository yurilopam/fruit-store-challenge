package services

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
	ttl    time.Duration
}

type RedisOption func(*RedisService)

func WithRedisTTLSeconds(sec int) RedisOption {
	return func(r *RedisService) { r.ttl = time.Duration(sec) * time.Second }
}

func NewRedis(addr string, db int, password string, opts ...RedisOption) *RedisService {
	rdb := redis.NewClient(&redis.Options{Addr: addr, DB: db, Password: password})
	rs := &RedisService{client: rdb, ttl: time.Minute}
	for _, opt := range opts {
		opt(rs)
	}
	return rs
}

func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisService) Set(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value, r.ttl).Err()
}

func (r *RedisService) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
