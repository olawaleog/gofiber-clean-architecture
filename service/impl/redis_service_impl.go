package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/go-redis/redis/v9"
)

// RedisService handles Redis operations for caching and pub/sub
type RedisService struct {
	Client *redis.Client
	Ctx    context.Context
}

// NewRedisService creates a new Redis service
func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		Client: client,
		Ctx:    context.Background(),
	}
}

// Set stores a value in Redis with an optional expiration time
func (s *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error marshaling data for Redis: %s", err.Error()))
		return err
	}

	err = s.Client.Set(s.Ctx, key, jsonData, expiration).Err()
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error setting data in Redis: %s", err.Error()))
		return err
	}

	common.Logger.Info(fmt.Sprintf("Data stored in Redis with key: %s", key))
	return nil
}

// Get retrieves a value from Redis and unmarshals it into the provided destination
func (s *RedisService) Get(key string, dest interface{}) error {
	val, err := s.Client.Get(s.Ctx, key).Result()
	if err == redis.Nil {
		common.Logger.Info(fmt.Sprintf("Key %s does not exist in Redis", key))
		return fmt.Errorf("key %s not found", key)
	} else if err != nil {
		common.Logger.Error(fmt.Sprintf("Error getting data from Redis: %s", err.Error()))
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error unmarshaling data from Redis: %s", err.Error()))
		return err
	}

	return nil
}

// Delete removes a key from Redis
func (s *RedisService) Delete(key string) error {
	err := s.Client.Del(s.Ctx, key).Err()
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error deleting key from Redis: %s", err.Error()))
		return err
	}

	common.Logger.Info(fmt.Sprintf("Key %s deleted from Redis", key))
	return nil
}

// PublishMessage publishes a message to a Redis channel
func (s *RedisService) PublishMessage(channel string, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error marshaling message for Redis pub/sub: %s", err.Error()))
		return err
	}

	err = s.Client.Publish(s.Ctx, channel, jsonData).Err()
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error publishing message to Redis channel: %s", err.Error()))
		return err
	}

	common.Logger.Info(fmt.Sprintf("Message published to Redis channel: %s", channel))
	return nil
}

// SubscribeToChannel subscribes to a Redis channel and processes messages with the handler
func (s *RedisService) SubscribeToChannel(channel string, handler func([]byte) error) error {
	pubsub := s.Client.Subscribe(s.Ctx, channel)
	defer pubsub.Close()

	// Confirm subscription
	_, err := pubsub.Receive(s.Ctx)
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error confirming subscription to Redis channel: %s", err.Error()))
		return err
	}

	// Process messages in a goroutine
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			err := handler([]byte(msg.Payload))
			if err != nil {
				common.Logger.Error(fmt.Sprintf("Error handling message from Redis channel: %s", err.Error()))
			}
		}
	}()

	common.Logger.Info(fmt.Sprintf("Subscribed to Redis channel: %s", channel))
	return nil
}

// GetKeys returns all keys matching the pattern
func (s *RedisService) GetKeys(pattern string) ([]string, error) {
	keys, err := s.Client.Keys(s.Ctx, pattern).Result()
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error getting keys matching pattern %s: %s", pattern, err.Error()))
		return nil, err
	}
	return keys, nil
}

// Exists checks if a key exists
func (s *RedisService) Exists(key string) (bool, error) {
	result, err := s.Client.Exists(s.Ctx, key).Result()
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error checking if key %s exists: %s", key, err.Error()))
		return false, err
	}
	return result > 0, nil
}

// SetWithTTL stores a value in Redis with a specific TTL
func (s *RedisService) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error marshaling data for Redis: %s", err.Error()))
		return err
	}

	err = s.Client.Set(s.Ctx, key, jsonData, ttl).Err()
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error setting data with TTL in Redis: %s", err.Error()))
		return err
	}

	common.Logger.Info(fmt.Sprintf("Data stored in Redis with key: %s and TTL: %v", key, ttl))
	return nil
}

// GetTTL returns the remaining TTL of a key
func (s *RedisService) GetTTL(key string) (time.Duration, error) {
	ttl, err := s.Client.TTL(s.Ctx, key).Result()
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error getting TTL for key %s: %s", key, err.Error()))
		return 0, err
	}
	return ttl, nil
}
