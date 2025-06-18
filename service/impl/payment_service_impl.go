package impl

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
)

type PaymentServiceImpl struct {
	repository.PaymentRepository
}

func NewPaymentService(paymentRepository *repository.PaymentRepository) service.PaymentService {
	return &PaymentServiceImpl{*paymentRepository}
}

func (p PaymentServiceImpl) GetPaymentMethods(ctx context.Context, userId float64) interface{} {
	var paymentMethodModels []model.PaymentMethodModel
	paymentMethods := p.PaymentRepository.ListPaymentMethods(ctx, userId)
	for _, paymentMethod := range paymentMethods {
		paymentMethodModels = append(paymentMethodModels, model.PaymentMethodModel{
			UniqueId: paymentMethod.UniqueId,
			Name:     paymentMethod.Name,
			Provider: paymentMethod.Provider,
			Scheme:   paymentMethod.Scheme,
			AuthCode: paymentMethod.AuthCode,
			Id:       paymentMethod.ID,
		})
	}
	return paymentMethodModels
}

func (p PaymentServiceImpl) InitiateMobileMoneyTransaction(ctx context.Context, data model.MobileMoneyRequestModel) interface{} {
	//response := p.PaymentRepository.InitialiseMobileMoneyTransaction(ctx, data)
	panic("implement me")
}
