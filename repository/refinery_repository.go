package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
)

type RefineryRepository interface {
	ListRefinery(ctx context.Context) ([]entity.Refinery, error)
	Create(ctx context.Context, refinery entity.Refinery) (entity.Refinery, error)
	Update(ctx context.Context, refinery entity.Refinery, id string) (entity.Refinery, error)
	GetRefineryDashboardData(ctx context.Context, u uint) (map[string]interface{}, error)
}
