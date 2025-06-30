package impl

import (
	"context"
	"errors"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
	"time"
)

func NewTransactionRepositoryImpl(DB *gorm.DB) repository.TransactionRepository {
	return &transactionRepositoryImpl{DB: DB}
}

type transactionRepositoryImpl struct {
	*gorm.DB
}

func (transactionRepository *transactionRepositoryImpl) GetAdminDashboardData(ctx context.Context) map[string]interface{} {
	var refineries []entity.Refinery
	var orders []entity.Order
	var users []entity.User
	var data map[string]interface{}
	var totalOrderCount int
	var totalRefineryCount int
	var customerCount int
	var totalOrderAmount float64
	err := transactionRepository.DB.WithContext(ctx).
		Preload("Transaction").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Where("tb_transactions.status = ?", "success").
		Find(&orders).Error
	exception.PanicLogging(err)
	for _, order := range orders {
		if order.Status == 0 {
			totalOrderCount++
		} else if order.Status == 3 {
			totalOrderCount++
		}
		totalOrderAmount += order.Transaction.Amount
	}
	err = transactionRepository.DB.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&refineries).Error
	exception.PanicLogging(err)
	for _, refinery := range refineries {
		if refinery.IsActive {
			totalRefineryCount++
		}
	}
	err = transactionRepository.DB.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&users).Error

	exception.PanicLogging(err)
	for _, user := range users {
		if user.IsActive && user.UserRole == "customer" {
			customerCount++
		}
	}

	data = map[string]interface{}{
		"totalOrderCount":    totalOrderCount,
		"totalRefineryCount": totalRefineryCount,
		"totalOrderAmount":   totalOrderAmount,
		"customerCount":      customerCount,
	}
	return data
}

func (transactionRepository *transactionRepositoryImpl) FindByReference(ctx context.Context, id string) (entity.Transaction, error) {
	var transaction entity.Transaction
	result := transactionRepository.DB.Where("id = ?", id).First(&transaction)
	if result.RowsAffected == 0 {
		return entity.Transaction{}, errors.New("transaction Not Found")
	}
	return transaction, nil
}

func (transactionRepository *transactionRepositoryImpl) Insert(ctx context.Context, transaction entity.Transaction) entity.Transaction {
	err := transactionRepository.DB.WithContext(ctx).Create(&transaction).Error
	exception.PanicLogging(err)
	return transaction
}

func (transactionRepository *transactionRepositoryImpl) Delete(ctx context.Context, transaction entity.Transaction) {
	transactionRepository.DB.WithContext(ctx).Delete(&transaction)
}

func (transactionRepository *transactionRepositoryImpl) FindById(ctx context.Context, id string) (entity.Transaction, error) {
	var transaction entity.Transaction
	result := transactionRepository.DB.WithContext(ctx).
		Table("tb_transaction").
		Select("tb_transaction.transaction_id, tb_transaction.total_price, tb_transaction_detail.transaction_detail_id, tb_transaction_detail.sub_total_price, tb_transaction_detail.price, tb_transaction_detail.quantity, tb_product.product_id, tb_product.name, tb_product.price, tb_product.quantity").
		Joins("join tb_transaction_detail on tb_transaction_detail.transaction_id = tb_transaction.transaction_id").
		Joins("join tb_product on tb_product.product_id = tb_transaction_detail.product_id").
		Preload("TransactionDetails").
		Preload("TransactionDetails.Product").
		Where("tb_transaction.transaction_id = ?", id).
		First(&transaction)
	if result.RowsAffected == 0 {
		return entity.Transaction{}, errors.New("transaction Not Found")
	}
	return transaction, nil
}

func (transactionRepository *transactionRepositoryImpl) FindAll(ctx context.Context) []entity.Transaction {
	var transactions []entity.Transaction
	transactionRepository.DB.WithContext(ctx).
		Order("created_at desc").
		Find(&transactions)
	return transactions
}

func (transactionRepository *transactionRepositoryImpl) Update(ctx context.Context, transaction entity.Transaction) error {
	err := transactionRepository.DB.WithContext(ctx).Save(&transaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (transactionRepository *transactionRepositoryImpl) GetRefineryDashboardData(ctx context.Context, u uint) (map[string]interface{}, error) {
	var orders []entity.Order
	var data map[string]interface{}
	var pendingRequestsCount int
	var processedRequestsCount int
	var revenue float64

	err := transactionRepository.DB.WithContext(ctx).
		Preload("Transaction").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Where("tb_transactions.status = ?", "success").
		Find(&orders).Error
	exception.PanicLogging(err)
	for _, order := range orders {
		if order.Status == 0 {
			pendingRequestsCount++
		} else if order.Status == 3 {
			processedRequestsCount++
		}
		revenue += order.Transaction.Amount
	}

	data = map[string]interface{}{
		"pendingRequests":   pendingRequestsCount,
		"processedRequests": processedRequestsCount,
		"revenue":           revenue,
	}
	return data, nil

}

// Add this method to your TransactionRepositoryImpl
func (repository *transactionRepositoryImpl) FindPendingTransactionsOlderThan(ctx context.Context, duration time.Duration) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	result := repository.DB.WithContext(ctx).
		Where("status = ? ", "Initiated").
		Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}
