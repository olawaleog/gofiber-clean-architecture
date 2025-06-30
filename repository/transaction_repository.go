package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"time"
)

type TransactionRepository interface {
	Insert(ctx context.Context, transaction entity.Transaction) entity.Transaction
	Delete(ctx context.Context, transaction entity.Transaction)
	FindById(ctx context.Context, id string) (entity.Transaction, error)
	FindAll(ctx context.Context) []entity.Transaction
	Update(ctx context.Context, transaction entity.Transaction) error
	FindByReference(ctx context.Context, id string) (entity.Transaction, error)
	GetRefineryDashboardData(ctx context.Context, u uint) (map[string]interface{}, error)
	GetAdminDashboardData(ctx context.Context) map[string]interface{}
	FindPendingTransactionsOlderThan(ctx context.Context, duration time.Duration) ([]entity.Transaction, error)
}
