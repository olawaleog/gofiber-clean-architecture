package service

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type PaymentConfigService interface {
	GetPaymentConfig(ctx context.Context, countryCode string) (*model.PaymentConfigModel, error)
	ListPaymentConfigs(ctx context.Context) ([]model.PaymentConfigModel, error)
	CreatePaymentConfig(ctx context.Context, config *model.PaymentConfigModel) error
	UpdatePaymentConfig(ctx context.Context, config *model.PaymentConfigModel) error
	DeletePaymentConfig(ctx context.Context, countryCode string) error
	RefreshCache(ctx context.Context) error
}
