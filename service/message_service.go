package service

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type MessageService interface {
	SendEmail(ctx context.Context, model model.EmailMessageModel)
	SendSMS(ctx context.Context, model model.SMSMessageModel) error
	FindMessageTemplateById(ctx context.Context, id int) model.MessageTemplateModel
	FindAllMessageTemplate(ctx context.Context) []model.MessageTemplateModel
	CreateMessageTemplate(ctx context.Context, model model.MessageTemplateModel)
	UpdateMessageTemplate(ctx context.Context, model model.MessageTemplateModel)
	GenerateOneTimePassword(ctx context.Context, uid uint) (entity.OneTimePassword, error)
	FindMessageTemplateByName(ctx context.Context, name string) model.MessageTemplateModel
	SendSMSDirect(data model.SMSMessageModel) error
}
