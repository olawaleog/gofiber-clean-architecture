package service

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type TruckService interface {
	ListAllTrucks(ctx context.Context) ([]model.TruckModel, error)
	CreateTruck(truck model.TruckModel) (model.TruckModel, error)
	UpdateTruck(truck model.TruckModel) (model.TruckModel, error)
	GetActiveTruck(ctx context.Context) model.TruckModel
}
