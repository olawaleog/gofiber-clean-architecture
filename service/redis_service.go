package service

import (
	"time"
)

type RedisService interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
}
