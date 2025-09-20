package service

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		Client: client,
	}
}

// Get retrieves a value from Redis by key
func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Set stores a value in Redis with a specified TTL
func (r *RedisService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a key from Redis
func (r *RedisService) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
