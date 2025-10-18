package impl

import (
	"context"
	"strconv"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
)

type TruckRepositoryImpl struct {
	*gorm.DB
}

func NewTruckRepositoryImpl(db *gorm.DB) repository.TruckRepository {
	return &TruckRepositoryImpl{db}

}

func (t TruckRepositoryImpl) ListTrucks(ctx context.Context) ([]model.TruckModel, error) {
	var trucks []entity.Truck
	err := t.DB.WithContext(ctx).
		Preload("User").
		Where("is_active = ?", true).
		Find(&trucks).Error
	exception.PanicLogging(err)
	var truckModels []model.TruckModel
	for _, truck := range trucks {
		truckModels = append(truckModels, model.TruckModel{
			Id:                    truck.ID,
			Capacity:              strconv.Itoa(truck.Capacity),
			PlateNumber:           truck.PlateNumber,
			IsActive:              truck.IsActive,
			ManufacturerModel:     truck.ManufacturerModel,
			YearOfManufacture:     strconv.Itoa(truck.YearOfManufacture),
			EngineNumber:          truck.EngineNumber,
			LicenceExpirationDate: truck.LicenceExpirationDate,
			UserId:                truck.UserId,
			User: model.UserModel{
				Id:           truck.User.ID,
				Username:     truck.User.Username,
				Role:         truck.User.UserRole,
				PhoneNumber:  truck.User.PhoneNumber,
				EmailAddress: truck.User.Email,
				LastName:     truck.User.LastName,
				FirstName:    truck.User.FirstName,
				IsActive:     truck.User.IsActive,
			},
		})
	}
	return truckModels, nil
}

func (t TruckRepositoryImpl) Create(truck entity.Truck) (entity.Truck, error) {
	err := t.DB.Create(&truck).Error
	exception.PanicLogging(err)
	return truck, nil
}

func (t TruckRepositoryImpl) UpdateTruck(truck model.TruckModel) (model.TruckModel, error) {
	var truckUpdate entity.Truck
	err := t.DB.Where("id = ?", truck.Id).
		Preload("User").
		Joins("JOIN tb_users ON tb_users.id = tb_trucks.user_id").
		First(&truckUpdate).Error
	exception.PanicLogging(err)
	capacity, err := strconv.Atoi(truck.Capacity)
	truckUpdate.Capacity = capacity
	truckUpdate.PlateNumber = truck.PlateNumber
	truckUpdate.UserId = truck.UserId
	err = t.DB.Save(&truckUpdate).Error
	exception.PanicLogging(err)
	return truck, nil
}

func (t TruckRepositoryImpl) GetActiveTruck(ctx context.Context) (entity.Truck, error) {
	var truck entity.Truck
	err := t.DB.WithContext(ctx).
		Where("tb_trucks.is_active = ?", true).
		Preload("User").
		Joins("JOIN tb_users ON tb_users.id = tb_trucks.user_id").
		First(&truck).Error
	return truck, err
}

// ListTrucksByCountryCode returns active trucks whose owner's area_code matches the provided countryCode.
func (t TruckRepositoryImpl) ListTrucksByCountryCode(ctx context.Context, countryCode string) ([]model.TruckModel, error) {
	var trucks []entity.Truck
	var err error
	if countryCode == "All" {
		err = t.DB.WithContext(ctx).
			Preload("User").
			Joins("JOIN tb_users ON tb_users.id = tb_trucks.user_id").
			Where("tb_trucks.is_active = ? ", true).
			Find(&trucks).Error
	} else {
		err = t.DB.WithContext(ctx).
			Preload("User").
			Joins("JOIN tb_users ON tb_users.id = tb_trucks.user_id").
			Where("tb_trucks.is_active = ? AND tb_users.country_code = ?", true, countryCode).
			Find(&trucks).Error
	}
	exception.PanicLogging(err)

	var truckModels []model.TruckModel
	for _, truck := range trucks {
		truckModels = append(truckModels, model.TruckModel{
			Id:                    truck.ID,
			Capacity:              strconv.Itoa(truck.Capacity),
			PlateNumber:           truck.PlateNumber,
			IsActive:              truck.IsActive,
			ManufacturerModel:     truck.ManufacturerModel,
			YearOfManufacture:     strconv.Itoa(truck.YearOfManufacture),
			EngineNumber:          truck.EngineNumber,
			LicenceExpirationDate: truck.LicenceExpirationDate,
			UserId:                truck.UserId,
			User: model.UserModel{
				Id:           truck.User.ID,
				Username:     truck.User.Username,
				Role:         truck.User.UserRole,
				PhoneNumber:  truck.User.PhoneNumber,
				EmailAddress: truck.User.Email,
				LastName:     truck.User.LastName,
				FirstName:    truck.User.FirstName,
				IsActive:     truck.User.IsActive,
			},
		})
	}
	return truckModels, nil
}
