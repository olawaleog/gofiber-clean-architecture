package service

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type RefineryService interface {
	GetRefinery(context context.Context, model model.GetRefineryModel) (model.RefineryCostModel, error)
	ListRefineries(ctx context.Context) ([]model.RefineryModel, error)
	CreateRefinery(ctx context.Context, refineryModel model.CreateRefineryModel) (model.RefineryModel, error)
	UpdateRefinery(ctx context.Context, refineryModel model.CreateRefineryModel, id string) (model.RefineryModel, error)
	GetRefineryDashboardData(ctx context.Context, u uint) (map[string]interface{}, error)
	ToggleRefineryStatus(ctx context.Context, statusModel model.ToggleRefineryStatusModel) (bool, error)
}
