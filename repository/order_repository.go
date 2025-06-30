package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"time"
)

type OrderRepository interface {
	FindByTransactionId(ctx context.Context, transactionId uint) (entity.Order, error)
	Insert(ctx context.Context, order entity.Order) entity.Order
	GetRefineryOrders(ctx context.Context, u uint) ([]entity.Order, error)
	FindById(ctx context.Context, id uint) (entity.Order, error)
	Update(ctx context.Context, order entity.Order) error
	FindDriverOrdersByUserId(ctx context.Context, id float64, stage uint) ([]entity.Order, error)
	GetUserOrders(ctx context.Context, u uint) ([]entity.Order, error)
	FindInitiatedOrders(ctx context.Context, duration time.Duration) ([]entity.Order, error)
}
