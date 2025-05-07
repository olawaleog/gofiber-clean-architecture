package impl

import (
	"context"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"gopkg.in/gomail.v2"
	"strconv"
)

func NewMessageServiceImpl(config configuration.Config, messageTemplateRepository repository.MessageTemplateRepository, client *service.HttpService) service.MessageService {
	return &messageServiceImpl{config: config, MessageTemplateRepository: messageTemplateRepository, HttpService: *client}
}

type messageServiceImpl struct {
	service.HttpService
	config configuration.Config
	repository.MessageTemplateRepository
}

func (m *messageServiceImpl) GenerateOneTimePassword(context context.Context, uid uint) (entity.OneTimePassword, error) {
	otp, err := m.MessageTemplateRepository.GenerateOneTImePassword(context, uid)
	exception.PanicLogging(err)
	return otp, nil
}

func (m *messageServiceImpl) SendSMS(ctx context.Context, data model.SMSMessageModel) {
	headers := make(map[string]interface{})
	headers["apiKey"] = m.config.Get("AFRICAS_TALKING_API_KEY")
	body := make(map[string]interface{})
	body["username"] = m.config.Get("AFRICAS_TALKING_USERNAME")
	body["enqueue"] = 1
	body["message"] = data.Message

	body["to"] = "+234" + data.PhoneNumber[1:]
	url := m.config.Get("AFRICAS_TALKING_BASE_URL") + "/version1/messaging"
	m.HttpService.PostMethod(ctx, url, "POST", &body, &headers, true)
}

func (m *messageServiceImpl) FindMessageTemplateById(ctx context.Context, id int) model.MessageTemplateModel {
	var messageTemplate entity.MessageTemplate
	messageTemplate, err := m.MessageTemplateRepository.FindById(ctx, id)
	exception.PanicLogging(err)
	result := model.MessageTemplateModel{
		Name:    messageTemplate.Name,
		Message: messageTemplate.Message,
		Subject: messageTemplate.Subject,
	}
	return result
}

func (m *messageServiceImpl) FindMessageTemplateByName(ctx context.Context, name string) model.MessageTemplateModel {
	messageTemplate, err := m.MessageTemplateRepository.FindByName(ctx, name)
	exception.PanicLogging(err)
	messageTemplateModel := model.MessageTemplateModel{
		Name:    messageTemplate.Name,
		Message: messageTemplate.Message,
		Subject: messageTemplate.Subject,
		IsEmail: messageTemplate.IsEmailMessage,
		IsSms:   messageTemplate.IsSMSMessage,
	}
	return messageTemplateModel
}

func (m *messageServiceImpl) FindAllMessageTemplate(ctx context.Context) []model.MessageTemplateModel {
	var messageTemplates []entity.MessageTemplate
	messageTemplates, err := m.MessageTemplateRepository.FindAll(ctx)
	exception.PanicLogging(err)
	var result []model.MessageTemplateModel
	for _, messageTemplate := range messageTemplates {
		result = append(result, model.MessageTemplateModel{
			Name:    messageTemplate.Name,
			Message: messageTemplate.Message,
			Subject: messageTemplate.Subject,
		})
	}
	return result
}

func (m *messageServiceImpl) CreateMessageTemplate(ctx context.Context, model model.MessageTemplateModel) {
	var messageTemplate entity.MessageTemplate
	messageTemplate.Name = model.Name
	messageTemplate.Message = model.Message
	messageTemplate.Subject = model.Subject
	err := m.MessageTemplateRepository.Create(ctx, model)
	exception.PanicLogging(err)
}

func (m *messageServiceImpl) UpdateMessageTemplate(ctx context.Context, template model.MessageTemplateModel) {
	err := m.MessageTemplateRepository.Update(ctx, template)
	exception.PanicLogging(err)
}

func (m *messageServiceImpl) SendEmail(context context.Context, model model.EmailMessageModel) {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", m.config.Get("MAIL_FROM"))
	message.SetHeader("To", model.To)
	message.SetHeader("Subject", model.Subject)

	// Set email body
	message.SetBody("text/html", model.Message)

	// Set up the SMTP
	port, _ := strconv.Atoi(m.config.Get("MAIL_PORT"))
	dialer := gomail.NewDialer(m.config.Get("MAIL_HOST"), port, m.config.Get("MAIL_USERNAME"), m.config.Get("MAIL_PASSWORD"))

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
	return
}
