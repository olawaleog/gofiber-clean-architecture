package impl

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
)

type PaymentMethodRepositoryImpl struct {
	db *gorm.DB
}

func (p PaymentMethodRepositoryImpl) GetAll(ctx context.Context) ([]entity.PaymentMethod, error) {
	//TODO implement me
	panic("implement me")
}

func (p PaymentMethodRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]entity.PaymentMethod, error) {
	//TODO implement me
	panic("implement me")
}

func (p PaymentMethodRepositoryImpl) GetByID(ctx context.Context, paymentMethodID string) (entity.PaymentMethod, error) {
	//TODO implement me
	var paymentMethod entity.PaymentMethod
	err := p.db.Where("id = ?", paymentMethodID).First(&paymentMethod).Error
	return paymentMethod, err
}

func (p PaymentMethodRepositoryImpl) Create(ctx context.Context, paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error) {
	var existingPaymentMethod entity.PaymentMethod
	err := p.db.Where("unique_id = ? AND user_id = ?", paymentMethod.UniqueId, paymentMethod.UserID).First(&existingPaymentMethod).Error
	if err == nil {
		return existingPaymentMethod, nil
	}
	err = p.db.Create(&paymentMethod).Error
	return paymentMethod, err
}

func (p PaymentMethodRepositoryImpl) Delete(ctx context.Context, id int) {
	//TODO implement me
	panic("implement me")
}

func NewPaymentMethodRepository(db *gorm.DB) repository.PaymentMethodRepository {
	return &PaymentMethodRepositoryImpl{db}
}
