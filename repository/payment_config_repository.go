package repository

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
)

type PaymentConfigRepository interface {
	GetPaymentConfig(ctx context.Context, providerCode string) (*entity.PaymentConfiguration, error)
	ListPaymentConfigs(ctx context.Context) ([]entity.PaymentConfiguration, error)
	SavePaymentConfig(ctx context.Context, config *entity.PaymentConfiguration) error
	UpdatePaymentConfig(ctx context.Context, config *entity.PaymentConfiguration) error
	DeletePaymentConfig(ctx context.Context, providerCode string) error
}
