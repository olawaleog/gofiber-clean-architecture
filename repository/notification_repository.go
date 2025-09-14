package repository

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification entity.Notification) (entity.Notification, error)
	FindByUserID(ctx context.Context, userID uint) ([]entity.Notification, error)
	FindAll(ctx context.Context, page, size int) ([]entity.Notification, error)
	CountAll() (int64, error)
}
