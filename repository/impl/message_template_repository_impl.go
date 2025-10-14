package impl

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
)

func NewMessageTemplateRepositoryImpl(DB *gorm.DB) repository.MessageTemplateRepository {
	return &messageTemplateRepositoryImpl{
		DB: DB,
	}
}

type messageTemplateRepositoryImpl struct {
	*gorm.DB
}

func (m messageTemplateRepositoryImpl) FindByName(ctx context.Context, name string) (entity.MessageTemplate, error) {
	var message entity.MessageTemplate
	err := m.DB.WithContext(ctx).Where("template_name = ?", name).First(&message).Error
	exception.PanicLogging(err)
	return message, nil
}

func (m messageTemplateRepositoryImpl) FindById(ctx context.Context, i int) (entity.MessageTemplate, error) {
	var message entity.MessageTemplate
	err := m.DB.WithContext(ctx).Where("id = ?", i).First(&message)
	exception.PanicLogging(err)
	return message, nil
}

func (m messageTemplateRepositoryImpl) FindAll(ctx context.Context) ([]entity.MessageTemplate, error) {
	var messages []entity.MessageTemplate
	err := m.DB.WithContext(ctx).Find(&messages)
	exception.PanicLogging(err)
	return messages, nil
}

func (m messageTemplateRepositoryImpl) Create(ctx context.Context, model model.MessageTemplateModel) error {
	message := entity.MessageTemplate{
		Subject: model.Subject,
		Message: model.Message,
		Name:    model.Name,
	}
	err := m.DB.WithContext(ctx).Create(&message)
	exception.PanicLogging(err)
	return nil
}

func (m messageTemplateRepositoryImpl) Update(ctx context.Context, templateModel model.MessageTemplateModel) error {
	var message entity.MessageTemplate
	err := m.DB.WithContext(ctx).Where("name = ?", templateModel.Name).First(&message)
	exception.PanicLogging(err)
	message.Subject = templateModel.Subject
	message.Message = templateModel.Message
	err = m.DB.WithContext(ctx).Save(&message)
	exception.PanicLogging(err)
	return nil
}

func (m messageTemplateRepositoryImpl) GenerateOneTImePassword(ctx context.Context, userId uint) (entity.OneTimePassword, error) {
	var otp entity.OneTimePassword
	code := generateOTP()
	otp.Code = code
	otp.UserId = userId
	otp.IsUsed = false
	otp.ExpiredAt = time.Now().Add(15 * time.Minute)
	err := m.DB.WithContext(ctx).Create(&otp).Error
	exception.PanicLogging(err)
	return otp, nil
}

func generateOTP() string {
	// Generate a random 6-digit OTP
	otp := fmt.Sprintf("%04d", rand.Intn(10000))
	return otp
}
