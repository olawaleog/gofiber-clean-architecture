package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/go-redis/redis/v9"
)

// RedisServiceImpl handles Redis operations for caching and pub/sub
type RedisServiceImpl struct {
	Client *redis.Client
	Ctx    context.Context
}

// NewRedisService creates a new Redis service
func NewRedisService(client *redis.Client) service.RedisService {
	return &RedisServiceImpl{
		Client: client,
		Ctx:    context.Background(),
	}
}

// Set stores a value in Redis with an optional expiration time
func (s *RedisServiceImpl) Set(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error marshaling data for Redis: %s", err.Error()))
		return err
	}

	err = s.Client.Set(s.Ctx, key, jsonData, expiration).Err()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error setting data in Redis: %s", err.Error()))
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Data stored in Redis with key: %s", key))
	return nil
}

// Get retrieves a value from Redis
func (s *RedisServiceImpl) Get(key string) (string, error) {
	val, err := s.Client.Get(s.Ctx, key).Result()
	if err == redis.Nil {
		logger.Logger.Info(fmt.Sprintf("Key %s does not exist in Redis", key))
		return "", fmt.Errorf("key %s not found", key)
	} else if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error getting data from Redis: %s", err.Error()))
		return "", err
	}

	return val, nil
}

// Delete removes a key from Redis
func (s *RedisServiceImpl) Delete(key string) error {
	err := s.Client.Del(s.Ctx, key).Err()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error deleting key from Redis: %s", err.Error()))
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Key %s deleted from Redis", key))
	return nil
}

// PublishMessage publishes a message to a Redis channel
func (s *RedisServiceImpl) PublishMessage(channel string, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error marshaling message for Redis pub/sub: %s", err.Error()))
		return err
	}

	err = s.Client.Publish(s.Ctx, channel, jsonData).Err()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error publishing message to Redis channel: %s", err.Error()))
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Message published to Redis channel: %s", channel))
	return nil
}

// SubscribeToChannel subscribes to a Redis channel and processes messages with the handler
func (s *RedisServiceImpl) SubscribeToChannel(channel string, handler func([]byte) error) error {
	pubsub := s.Client.Subscribe(s.Ctx, channel)
	defer pubsub.Close()

	// Confirm subscription
	_, err := pubsub.Receive(s.Ctx)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error confirming subscription to Redis channel: %s", err.Error()))
		return err
	}

	// Process messages in a goroutine
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			err := handler([]byte(msg.Payload))
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("Error handling message from Redis channel: %s", err.Error()))
			}
		}
	}()

	logger.Logger.Info(fmt.Sprintf("Subscribed to Redis channel: %s", channel))
	return nil
}

// GetKeys returns all keys matching the pattern
func (s *RedisServiceImpl) GetKeys(pattern string) ([]string, error) {
	keys, err := s.Client.Keys(s.Ctx, pattern).Result()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error getting keys matching pattern %s: %s", pattern, err.Error()))
		return nil, err
	}
	return keys, nil
}

// Exists checks if a key exists
func (s *RedisServiceImpl) Exists(key string) (bool, error) {
	result, err := s.Client.Exists(s.Ctx, key).Result()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error checking if key %s exists: %s", key, err.Error()))
		return false, err
	}
	return result > 0, nil
}

// SetWithTTL stores a value in Redis with a specific TTL
func (s *RedisServiceImpl) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error marshaling data for Redis: %s", err.Error()))
		return err
	}

	err = s.Client.Set(s.Ctx, key, jsonData, ttl).Err()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error setting data with TTL in Redis: %s", err.Error()))
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Data stored in Redis with key: %s and TTL: %v", key, ttl))
	return nil
}

// GetTTL returns the remaining TTL of a key
func (s *RedisServiceImpl) GetTTL(key string) (time.Duration, error) {
	ttl, err := s.Client.TTL(s.Ctx, key).Result()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error getting TTL for key %s: %s", key, err.Error()))
		return 0, err
	}
	return ttl, nil
}
