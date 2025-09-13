package repository

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification entity.Notification) (entity.Notification, error)
	FindByUserID(ctx context.Context, userID uint) ([]entity.Notification, error)
	FindAll(ctx context.Context, page, size int) ([]entity.Notification, error)
	CountAll() (int64, error)
}

type notificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{
		DB: db,
	}
}

func (r *notificationRepository) Create(ctx context.Context, notification entity.Notification) (entity.Notification, error) {
	err := r.DB.WithContext(ctx).Create(&notification).Error
	return notification, err
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID uint) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindAll(ctx context.Context, page, size int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	offset := (page - 1) * size
	err := r.DB.WithContext(ctx).Offset(offset).Limit(size).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) CountAll() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Notification{}).Count(&count).Error
	return count, err
}
