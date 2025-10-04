package service

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type SettingService interface {
	Create(ctx context.Context, request model.SettingModel) (entity.Setting, error)
	Update(ctx context.Context, key string, request model.SettingModel) (entity.Setting, error)
	List(ctx context.Context) ([]model.SettingResponseModel, error)
	FindByKey(ctx context.Context, key string) (model.SettingResponseModel, error)
}
