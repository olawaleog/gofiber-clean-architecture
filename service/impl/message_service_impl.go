package impl

import (
	"context"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"gopkg.in/gomail.v2"
	"strconv"
)

func NewMessageServiceImpl(config configuration.Config, messageTemplateRepository repository.MessageTemplateRepository, httpService *service.HttpService, rabbitMQService *RabbitMQService) service.MessageService {
	return &messageServiceImpl{
		config:                    config,
		MessageTemplateRepository: messageTemplateRepository,
		HttpService:               *httpService,
		rabbitMQService:           rabbitMQService,
	}
}

type messageServiceImpl struct {
	service.HttpService
	config configuration.Config
	repository.MessageTemplateRepository
	rabbitMQService *RabbitMQService
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
	if data.PhoneNumber[0] == '0' {
		body["to"] = data.CountryCode + data.PhoneNumber[1:]
	} else {
		body["to"] = data.PhoneNumber
	}
	body["from"] = "AquaWizz"

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
	// Create a queued email message
	queuedEmail := model.QueuedEmailMessage{
		To:      model.To,
		Subject: model.Subject,
		Message: model.Message,
		From:    m.config.Get("MAIL_FROM"),
	}

	// Publish to RabbitMQ
	err := m.rabbitMQService.PublishMessage("email.send", queuedEmail)
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Failed to queue email to %s: %s", model.To, err.Error()))
		// Fall back to synchronous sending if publishing fails
		m.sendEmailDirect(queuedEmail)
	} else {
		common.Logger.Info(fmt.Sprintf("Email to %s queued successfully", model.To))
	}
	return
}

// sendEmailDirect sends an email directly as a fallback mechanism
func (m *messageServiceImpl) sendEmailDirect(email model.QueuedEmailMessage) {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", email.From)
	message.SetHeader("To", email.To)
	message.SetHeader("Subject", email.Subject)

	// Set email body
	message.SetBody("text/html", email.Message)

	// Set up the SMTP
	port, _ := strconv.Atoi(m.config.Get("MAIL_PORT"))
	dialer := gomail.NewDialer(m.config.Get("MAIL_HOST"), port, m.config.Get("MAIL_USERNAME"), m.config.Get("MAIL_PASSWORD"))

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		common.Logger.Error(fmt.Sprintf("Error sending email directly: %s", err.Error()))
	} else {
		common.Logger.Info("Email sent directly successfully!")
	}
}
