package repository

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
)

type PaymentConfigurationRepository interface {
	GetPaymentConfiguration(ctx context.Context) (*entity.PaymentConfiguration, error)
}
