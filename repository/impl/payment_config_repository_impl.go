package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type PaymentConfigRepositoryImpl struct {
	DB    *gorm.DB
	Redis *redis.Client
}

const (
	paymentConfigPrefix = "payment_config"
	defaultCacheExpiry  = 60 * time.Minute
)

func NewPaymentConfigRepository(db *gorm.DB, redis *redis.Client) repository.PaymentConfigRepository {
	return &PaymentConfigRepositoryImpl{
		DB:    db,
		Redis: redis,
	}
}

func (r *PaymentConfigRepositoryImpl) GetPaymentConfig(ctx context.Context, countryCode string) (*entity.PaymentConfiguration, error) {
	// Try to get from Redis first
	cacheKey := fmt.Sprintf("%s_%s", paymentConfigPrefix, countryCode)
	cachedData, err := r.Redis.Get(ctx, cacheKey).Bytes()

	if err == nil {
		var config entity.PaymentConfiguration
		if err := json.Unmarshal(cachedData, &config); err == nil {
			return &config, nil
		}
	}

	// If not in cache or error, get from DB
	var config entity.PaymentConfiguration
	if err := r.DB.WithContext(ctx).Where("country_code = ?", countryCode).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NotFoundError{Message: "Payment config not found"}
		}
		return nil, err
	}

	// Store in cache with default expiry
	if configBytes, err := json.Marshal(config); err == nil {
		r.Redis.Set(ctx, cacheKey, configBytes, defaultCacheExpiry)
	}

	return &config, nil
}

func (r *PaymentConfigRepositoryImpl) ListPaymentConfigs(ctx context.Context) ([]entity.PaymentConfiguration, error) {
	var configs []entity.PaymentConfiguration
	if err := r.DB.WithContext(ctx).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *PaymentConfigRepositoryImpl) SavePaymentConfig(ctx context.Context, config *entity.PaymentConfiguration) error {
	if err := r.DB.WithContext(ctx).Create(config).Error; err != nil {
		return err
	}

	// Cache the new config with default expiry
	cacheKey := fmt.Sprintf("%s_%s", paymentConfigPrefix, config.CountryCode)
	if configBytes, err := json.Marshal(config); err == nil {
		r.Redis.Set(ctx, cacheKey, configBytes, defaultCacheExpiry)
	}

	return nil
}

func (r *PaymentConfigRepositoryImpl) UpdatePaymentConfig(ctx context.Context, config *entity.PaymentConfiguration) error {
	var existing entity.PaymentConfiguration
	if err := r.DB.WithContext(ctx).Where("country_code = ?", config.CountryCode).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return exception.NotFoundError{Message: "Payment config not found"}
		}
		return err
	}

	// Update fields
	existing.PublicKey = config.PublicKey
	existing.SecretKey = config.SecretKey
	existing.UpdatedAt = time.Now()
	if err := r.DB.WithContext(ctx).Updates(existing).Error; err != nil {
		return err
	}

	// Update cache with default expiry
	cacheKey := fmt.Sprintf("%s_%s", paymentConfigPrefix, config.CountryCode)
	if configBytes, err := json.Marshal(config); err == nil {
		r.Redis.Set(ctx, cacheKey, configBytes, defaultCacheExpiry)
	}

	return nil
}

func (r *PaymentConfigRepositoryImpl) DeletePaymentConfig(ctx context.Context, countryCode string) error {
	if err := r.DB.WithContext(ctx).Where("country_code = ?", countryCode).Delete(&entity.PaymentConfiguration{}).Error; err != nil {
		return err
	}

	// Remove from cache
	cacheKey := fmt.Sprintf("%s_%s", paymentConfigPrefix, countryCode)
	r.Redis.Del(ctx, cacheKey)

	return nil
}
