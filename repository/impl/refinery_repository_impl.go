package impl

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
)

type RefineryRepositoryImpl struct {
	*gorm.DB
}

func NewRefineryRepositoryImpl(db *gorm.DB) repository.RefineryRepository {
	return &RefineryRepositoryImpl{db}
}

func (r RefineryRepositoryImpl) ListRefinery(ctx context.Context) ([]entity.Refinery, error) {
	var refinerys []entity.Refinery
	err := r.DB.WithContext(ctx).
		Order("created_at desc").
		Find(&refinerys).Error
	exception.PanicLogging(err)
	return refinerys, nil
}

func (r RefineryRepositoryImpl) Create(ctx context.Context, refinery entity.Refinery) (entity.Refinery, error) {
	err := r.DB.WithContext(ctx).Create(&refinery).Error
	if err != nil {
		return entity.Refinery{}, err
	}
	return refinery, nil
}
func (r RefineryRepositoryImpl) Update(ctx context.Context, refinery entity.Refinery, id string) (entity.Refinery, error) {
	err := r.DB.WithContext(ctx).Where("id = ?", id).Save(&refinery).Error
	if err != nil {
		return entity.Refinery{}, err
	}
	return refinery, nil
}

func (r RefineryRepositoryImpl) GetRefineryDashboardData(ctx context.Context, u uint) (map[string]interface{}, error) {
	var orders []entity.Order
	var data map[string]interface{}
	err := r.DB.WithContext(ctx).
		Preload("Transaction").
		Joins("JOIN tb_transactions ON tb_transactions.id = tb_orders.transaction_id").
		Where("tb_transactions.status = ?", "success").
		Find(&orders).Error

	exception.PanicLogging(err)
	count := len(orders)
	data["count"] = count
	return data, nil

}

func (r RefineryRepositoryImpl) FindById(ctx context.Context, id any) entity.Refinery {
	var refinery entity.Refinery
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&refinery).Error
	if err != nil {
		return entity.Refinery{}
	}
	return refinery
}
