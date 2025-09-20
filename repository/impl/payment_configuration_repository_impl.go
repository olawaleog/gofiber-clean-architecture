package impl

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"gorm.io/gorm"
)

type PaymentConfigurationRepositoryImpl struct {
	DB *gorm.DB
}

func NewPaymentConfigurationRepository(db *gorm.DB) *PaymentConfigurationRepositoryImpl {
	return &PaymentConfigurationRepositoryImpl{
		DB: db,
	}
}

func (r *PaymentConfigurationRepositoryImpl) GetPaymentConfiguration(ctx context.Context) (*entity.PaymentConfiguration, error) {
	var config entity.PaymentConfiguration
	if err := r.DB.First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}
