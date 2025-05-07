package impl

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &OrderRepositoryImpl{DB: *db}
}

func (o OrderRepositoryImpl) FindByTransactionId(ctx context.Context, transactionId uint) (entity.Order, error) {
	var order entity.Order
	err := o.DB.WithContext(ctx).Where("transaction_id = ?", transactionId).First(&order).Error
	if err != nil {
		return order, err
	}
	return order, nil
}

func (o OrderRepositoryImpl) Insert(ctx context.Context, order entity.Order) entity.Order {
	err := o.DB.WithContext(ctx).Create(&order).Error
	if err != nil {
		return order
	}
	return order
}

func (o OrderRepositoryImpl) FindById(ctx context.Context, id uint) (entity.Order, error) {
	var order entity.Order
	err := o.DB.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return order, err
	}
	return order, nil
}
func (o OrderRepositoryImpl) Update(ctx context.Context, order entity.Order) error {
	err := o.DB.WithContext(ctx).Save(&order).Error
	if err != nil {
		return err
	}
	return nil
}

func (o OrderRepositoryImpl) GetRefineryOrders(ctx context.Context, u uint) ([]entity.Order, error) {
	var order []entity.Order
	//today := time.Now().Format("2006-01-02") // Format the current date as YYYY-MM-DD

	err := o.DB.WithContext(ctx).
		Preload("Transaction").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Where("refinery_id = ? ", u).Find(&order).Error
	//Where("refinery_id = ? AND DATE(created_at) = ?", u, today).First(&order).Error
	if err != nil {
		return order, err
	}
	return order, nil
}

func (o OrderRepositoryImpl) FindDriverOrdersByUserId(ctx context.Context, id float64, stage uint) ([]entity.Order, error) {
	var truck entity.Truck
	err := o.DB.WithContext(ctx).
		Where("user_id = ?", id).First(&truck).Error
	exception.PanicLogging(err)
	var orders []entity.Order
	err = o.DB.WithContext(ctx).
		Preload("Transaction").
		Preload("Refinery").
		Preload("Transaction.User").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Where("truck_id = ? AND (tb_orders.status = ? or tb_orders.status = ?)", truck.ID, stage, stage+1).Find(&orders).Error
	exception.PanicLogging(err)
	return orders, nil

}
