package impl

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

// GetTransactionsPaginated retrieves a paginated list of transactions
func (t *transactionServiceImpl) GetTransactionsPaginated(ctx context.Context, page, limit int) ([]model.TransactionModel, int64) {
	var list []model.TransactionModel
	transactions, totalCount := t.TransactionRepository.FindAllPaginated(ctx, page, limit)

	for _, transaction := range transactions {
		list = append(list, model.TransactionModel{
			ID:          transaction.ID,
			Email:       transaction.Email,
			PhoneNumber: transaction.PhoneNumber,
			Amount:      transaction.Amount,
			Currency:    transaction.Currency,
			PaymentID:   transaction.PaymentID,
			Provider:    transaction.Provider,
			PaymentType: transaction.PaymentType,
			Status:      transaction.Status,
			Reference:   transaction.Reference,
			DeliveryFee: transaction.DeliveryFee,
			WaterCost:   transaction.WaterCost,
			CreatedAt:   transaction.CreatedAt,
		})
	}

	return list, totalCount
}
