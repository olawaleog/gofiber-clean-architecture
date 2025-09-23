package service

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type TransactionService interface {
	InitiateMobileMoneyTransaction(ctx context.Context, request model.MobileMoneyRequestModel) interface{}
	PaymentStatus(ctx context.Context, id string) map[string]interface{}
	GetRefineryDashboardData(ctx context.Context, u uint) (map[string]interface{}, error)
	GetRefineryOrders(ctx context.Context, u uint) ([]model.OrderModel, error)
	ApproveOrRejectOrder(ctx context.Context, orderModel model.ApproveOrRejectOrderModel) (interface{}, error)
	GetDriverPendingOrder(ctx context.Context, userId float64, stage uint) []model.OrderModel
	GetCustomerPendingOrder(ctx context.Context, userId float64, stage uint) []model.OrderModel
	GetTransactions(ctx context.Context) []model.TransactionModel
	GetTransactionsPaginated(ctx context.Context, page, limit int) ([]model.TransactionModel, int64)
	GetTransactionsByCountryCode(ctx context.Context, countryCode string, page, limit int) ([]model.TransactionModel, int64)
	GetAdminDashboardData(ctx context.Context) (map[string]interface{}, error)
	GetCustomerOrders(ctx context.Context, u uint) ([]model.OrderModel, error)
	FindById(ctx context.Context, id uint) (model.OrderModel, error)
	ProcessRecurringPayment(ctx context.Context, requestModel model.MobileMoneyRequestModel) interface{}
	ProcessPendingTransactions(ctx context.Context, truckId model.TruckModel) error
	MarkOrderReadyForDelivery(id string) error
	GetDriverCompletedOrder(ctx context.Context, id float64, i uint) []model.OrderModel
	CloseOrder(id string) error
	SubmitRating(ctx context.Context, ratingModel model.RatingModel) error
}
