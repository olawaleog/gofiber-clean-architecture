package impl

import (
	"context"
	"time"

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

func (o OrderRepositoryImpl) GetUserOrders(ctx context.Context, u uint) ([]entity.Order, error) {
	var orders []entity.Order
	err := o.DB.WithContext(ctx).
		Preload("Transaction").
		Preload("Refinery").
		Preload("Transaction.User").
		Preload("Transaction.Address").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Joins("JOIN tb_addresses ON tb_addresses.id = tb_transactions.address_id").
		Where("tb_users.id = ?", u).
		Order("tb_orders.id desc").
		Find(&orders).Error
	if err != nil {
		return orders, err
	}
	return orders, nil
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
	err := o.DB.WithContext(ctx).
		Preload("Transaction").
		Preload("Refinery").
		Preload("Transaction.User").
		Preload("Transaction.Address").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Where("tb_orders.id = ?", id).First(&order).Error
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

func (o OrderRepositoryImpl) GetRefineryOrders(ctx context.Context, u uint, countryCode string) ([]entity.Order, error) {
	var order []entity.Order
	//today := time.Now().Format("2006-01-02") // Format the current date as YYYY-MM-DD

	query := o.DB.WithContext(ctx).
		Preload("Transaction").
		Preload("Refinery").
		Preload("Transaction.User").
		Preload("Truck").
		Preload("Truck.User").
		Preload("Transaction.Address").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Joins("LEFT JOIN tb_trucks ON tb_trucks.id = tb_orders.truck_id").
		Joins(" LEFT JOIN tb_users AS truck_users ON truck_users.id = tb_trucks.user_id").
		Order("tb_orders.id desc").
		Order("tb_orders.status ")

	// Filter by country code if provided
	if countryCode != "" {
		query = query.Where("tb_transactions.country_code = ?", countryCode)
	}

	err := query.Order("tb_orders.id desc").Find(&order).Error

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
		Preload("Transaction.Address").
		Preload("Refinery").
		Preload("Transaction.User").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Where("truck_id = ? AND (tb_orders.status >= ?  AND tb_orders.status < ?)", truck.ID, 1, 4).
		Order("tb_orders.id desc").
		Find(&orders).Error
	exception.PanicLogging(err)
	return orders, nil

}
func (o OrderRepositoryImpl) FindCustomerOrdersByUserId(ctx context.Context, id float64, stage uint) ([]entity.Order, error) {

	var orders []entity.Order
	err := o.DB.WithContext(ctx).
		Preload("Transaction").
		Preload("Refinery").
		Preload("Transaction.User").
		Preload("Transaction.Address").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Where("tb_transactions.user_id = ? AND (tb_orders.status >= ?  AND tb_orders.status < ?)", id, 0, 4).
		Order("tb_orders.id desc").
		Find(&orders).Error
	exception.PanicLogging(err)
	return orders, nil

}

func (o OrderRepositoryImpl) FindCompletedDriverOrdersByUserId(ctx context.Context, id float64, stage uint) ([]entity.Order, error) {
	var truck entity.Truck
	err := o.DB.WithContext(ctx).
		Where("user_id = ?", id).First(&truck).Error
	exception.PanicLogging(err)
	var orders []entity.Order
	err = o.DB.WithContext(ctx).
		Preload("Transaction").
		Preload("Transaction.Address").
		Preload("Refinery").
		Preload("Transaction.User").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Where("truck_id = ? AND (tb_orders.status >= ? )", truck.ID, 4).
		Order("tb_orders.id desc").
		Find(&orders).Error
	exception.PanicLogging(err)
	return orders, nil

}

// Add this method to your TransactionRepositoryImpl
func (o *OrderRepositoryImpl) FindInitiatedOrders(ctx context.Context, duration time.Duration) ([]entity.Order, error) {
	var orders []entity.Order
	result := o.DB.WithContext(ctx).
		Preload("Transaction").
		Preload("Transaction.User").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Where("tb_orders.status = ? ", 0).
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (o OrderRepositoryImpl) MarkOrderReadyForDelivery(id string) (entity.Order, error) {
	var order entity.Order
	err := o.DB.WithContext(context.Background()).
		Preload("Transaction").
		Preload("Refinery").
		Preload("Transaction.User").
		//Preload("Truck").
		//Preload("Truck.User").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		//Joins("JOIN tb_trucks ON tb_trucks.id = tb_orders.truck_id").
		//Joins("JOIN tb_users AS truck_users ON truck_users.id = tb_trucks.user_id").
		Order("tb_orders.id desc").
		Where("tb_orders.id = ?", id).
		First(&order).Error
	if err != nil {
		return entity.Order{}, err
	}
	order.Status = 3
	err = o.DB.Save(&order).Error
	return order, err
}

func (o OrderRepositoryImpl) CloseOrder(id string) (entity.Order, error) {
	var order entity.Order
	err := o.DB.WithContext(context.Background()).
		Preload("Transaction").
		Preload("Refinery").
		Preload("Transaction.User").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Joins("JOIN tb_refineries ON tb_refineries.id = tb_orders.refinery_id").
		Joins("JOIN tb_users ON tb_users.id = tb_transactions.user_id ").
		Order("tb_orders.id desc").
		Where("tb_orders.id = ?", id).
		First(&order).Error
	if err != nil {
		return entity.Order{}, err
	}
	order.Status = 4
	err = o.DB.Save(&order).Error
	return order, err
}
