package service

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type PaymentService interface {
	GetPaymentMethods(context context.Context, userId float64) interface{}
	InitiateMobileMoneyTransaction(ctx context.Context, data model.MobileMoneyRequestModel) interface{}
}
