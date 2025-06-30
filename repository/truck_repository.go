package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type TruckRepository interface {
	ListTrucks(ctx context.Context) ([]model.TruckModel, error)
	Create(truck entity.Truck) (entity.Truck, error)
	UpdateTruck(truck model.TruckModel) (model.TruckModel, error)
	GetActiveTruck(ctx context.Context) (entity.Truck, error)
}
