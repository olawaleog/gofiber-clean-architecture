package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type SettingRepository interface {
	Create(ctx context.Context, setting model.SettingModel) (entity.Setting, error)
	Update(ctx context.Context, key string, setting model.SettingModel) (entity.Setting, error)
	List(ctx context.Context) ([]entity.Setting, error)
	FindByKey(ctx context.Context, key string) (entity.Setting, error)
}
