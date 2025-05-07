package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type PaymentRepository interface {
	ListPaymentMethods(ctx context.Context, id float64) []entity.PaymentMethod
	InitialiseMobileMoneyTransaction(ctx context.Context, data model.MobileMoneyRequestModel) interface{}
}
