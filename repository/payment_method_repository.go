package repository

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
)

type PaymentMethodRepository interface {
	GetAll(ctx context.Context) ([]entity.PaymentMethod, error)
	GetByUserID(ctx context.Context, userID string) ([]entity.PaymentMethod, error)
	GetByID(ctx context.Context, paymentMethodID string) (entity.PaymentMethod, error)
	Create(ctx context.Context, paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error)
	Delete(ctx context.Context, id int)
}
