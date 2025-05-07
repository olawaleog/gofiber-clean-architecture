package impl

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
)

type PaymentRepositoryimpl struct {
	*gorm.DB
}

func NewPaymentRepository(db *gorm.DB) repository.PaymentRepository {
	return &PaymentRepositoryimpl{DB: db}
}

func (p PaymentRepositoryimpl) ListPaymentMethods(ctx context.Context, id float64) []entity.PaymentMethod {
	var paymentMethods []entity.PaymentMethod
	err := p.DB.WithContext(ctx).Where("user_id = ?", id).Find(&paymentMethods).Error
	exception.PanicLogging(err)
	return paymentMethods
}

func (p PaymentRepositoryimpl) InitialiseMobileMoneyTransaction(ctx context.Context, data model.MobileMoneyRequestModel) interface{} {
	//TODO implement me
	panic("implement me")
}
