package impl

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"strconv"
)

type TruckServiceImpl struct {
	repository.TruckRepository
	service.UserService
	service.MessageService
}

func NewTruckServiceImpl(r *repository.TruckRepository, u *service.UserService, m *service.MessageService) service.TruckService {
	return &TruckServiceImpl{TruckRepository: *r, UserService: *u, MessageService: *m}
}

func (t TruckServiceImpl) ListAllTrucks(c context.Context) ([]model.TruckModel, error) {
	truck, err := t.TruckRepository.ListTrucks(c)
	exception.PanicLogging(err)
	return truck, nil
}

func (t TruckServiceImpl) CreateTruck(truckModel model.TruckModel) (model.TruckModel, error) {
	yearOfmanufacture, err := strconv.Atoi(truckModel.YearOfManufacture)
	capacity, err := strconv.Atoi(truckModel.Capacity)
	exception.PanicLogging(err)

	password, err := common.GeneratePassword(8)
	exception.PanicLogging(err)
	user := model.UserModel{
		Username:     truckModel.Phone,
		Password:     password,
		Role:         common.TRUCK_DRIVER_ROLE,
		EmailAddress: truckModel.Email,
		FirstName:    truckModel.FirstName,
		LastName:     truckModel.LastName,
		PhoneNumber:  truckModel.Phone,
		IsActive:     false,
		CountryCode:  "+233",
	}

	userResult := t.UserService.Register(context.TODO(), user)
	exception.PanicLogging(err)

	truckEntity := entity.Truck{
		ManufacturerModel:     truckModel.ManufacturerModel,
		YearOfManufacture:     yearOfmanufacture,
		PlateNumber:           truckModel.PlateNumber,
		Capacity:              capacity,
		EngineNumber:          truckModel.EngineNumber,
		IsActive:              true,
		UserId:                userResult.ID,
		LicenceExpirationDate: truckModel.LicenceExpirationDate,
	}
	truck, err := t.TruckRepository.Create(truckEntity)
	exception.PanicLogging(err)
	truckModel.Id = truck.ID

	smsModel := model.SMSMessageModel{
		Message:     "Hello " + truckModel.FirstName + "," + "Your Truck has been registered successfully.Your password is " + user.Password + ":\n Download  the app to complete your registration\n",
		PhoneNumber: user.PhoneNumber,
		CountryCode: userResult.CountryCode,
	}
	t.MessageService.SendSMS(context.TODO(), smsModel)

	return truckModel, nil
}

func (t TruckServiceImpl) UpdateTruck(truck model.TruckModel) (model.TruckModel, error) {
	truck, err := t.TruckRepository.UpdateTruck(truck)
	exception.PanicLogging(err)
	return truck, nil
}

func (t TruckServiceImpl) GetActiveTruck(ctx context.Context) model.TruckModel {
	truck, err := t.TruckRepository.GetActiveTruck(ctx)
	if err != nil {
		return model.TruckModel{}
	}
	return model.TruckModel{
		Id: truck.ID,
	}

}
