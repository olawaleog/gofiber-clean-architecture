package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type MessageTemplateRepository interface {
	FindById(context.Context, int) (entity.MessageTemplate, error)
	FindAll(context.Context) ([]entity.MessageTemplate, error)
	Create(context.Context, model.MessageTemplateModel) error
	Update(context.Context, model.MessageTemplateModel) error
	GenerateOneTImePassword(context.Context, uint) (entity.OneTimePassword, error)
	FindByName(ctx context.Context, name string) (entity.MessageTemplate, error)
}
