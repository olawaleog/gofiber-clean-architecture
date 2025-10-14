package impl

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"gopkg.in/gomail.v2"
)

func NewMessageServiceImpl(config configuration.Config, messageTemplateRepository repository.MessageTemplateRepository, httpService *service.HttpService, rabbitMQService *RabbitMQService, notificationRepository repository.NotificationRepository) service.MessageService {
	return &messageServiceImpl{
		config:                    config,
		MessageTemplateRepository: messageTemplateRepository,
		HttpService:               *httpService,
		rabbitMQService:           rabbitMQService,
		notificationRepository:    notificationRepository,
	}
}

type messageServiceImpl struct {
	service.HttpService
	config configuration.Config
	repository.MessageTemplateRepository
	rabbitMQService        *RabbitMQService
	notificationRepository repository.NotificationRepository
}

func (m *messageServiceImpl) GenerateOneTimePassword(context context.Context, uid uint) (entity.OneTimePassword, error) {
	otp, err := m.MessageTemplateRepository.GenerateOneTImePassword(context, uid)
	exception.PanicLogging(err)
	return otp, nil
}

func (m *messageServiceImpl) SendSMS(ctx context.Context, data model.SMSMessageModel) error {
	// Create a queued SMS message
	queuedSMS := model.QueuedSMSMessage{
		PhoneNumber: data.PhoneNumber,
		CountryCode: data.CountryCode,
		Message:     data.Message,
	}

	// Try to publish to RabbitMQ
	err := m.rabbitMQService.PublishMessage("sms.send", queuedSMS)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to queue SMS to %s: %s", data.PhoneNumber, err.Error()))
		// Fall back to synchronous sending if publishing fails
		m.SendSMSDirect(data)
		return err
	} else {
		logger.Logger.Info(fmt.Sprintf("SMS to %s queued successfully", data.PhoneNumber))
	}

	// Save notification record
	notification := entity.Notification{
		Type:      entity.SMS,
		Recipient: data.CountryCode + data.PhoneNumber,
		Content:   data.Message,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	_, err = m.notificationRepository.Create(ctx, notification)
	if err != nil {

		logger.Logger.Error(fmt.Sprintf("Failed to save SMS notification: %s", err.Error()))
		return err
	}
	return nil
}

// sendSMSDirect sends an SMS directly as a fallback mechanism
func (m *messageServiceImpl) SendSMSDirect(data model.SMSMessageModel) error {
	headers := make(map[string]interface{})
	headers["apiKey"] = m.config.Get("AFRICAS_TALKING_API_KEY")

	body := make(map[string]interface{})
	body["username"] = m.config.Get("AFRICAS_TALKING_USERNAME")
	body["enqueue"] = 1
	body["message"] = data.Message

	phoneNumber := data.PhoneNumber
	if phoneNumber == "" {
		return errors.New("invalid phone number")
	}
	if phoneNumber[0] == '0' {
		phoneNumber = data.CountryCode + phoneNumber[1:]
	} else if data.CountryCode != "" && phoneNumber[0] != '+' {
		phoneNumber = data.CountryCode + phoneNumber
	}

	body["to"] = phoneNumber
	body["from"] = "AquaWizz"

	url := m.config.Get("AFRICAS_TALKING_BASE_URL") + "/version1/messaging"
	_, err := m.HttpService.PostMethod(context.Background(), url, "POST", &body, &headers, true)

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error sending SMS directly: %s", err.Error()))
		return err
	} else {
		logger.Logger.Info(fmt.Sprintf("SMS sent directly to %s successfully!", phoneNumber))
	}
	return nil
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

func (m *messageServiceImpl) SendEmail(ctx context.Context, emailModel model.EmailMessageModel) {
	// Create a queued email message
	queuedEmail := model.QueuedEmailMessage{
		To:      emailModel.To,
		Subject: emailModel.Subject,
		Message: emailModel.Message,
		From:    m.config.Get("MAIL_FROM"),
	}

	// Publish to RabbitMQ
	err := m.rabbitMQService.PublishMessage("email.send", queuedEmail)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to queue email to %s: %s", emailModel.To, err.Error()))
		// Fall back to synchronous sending if publishing fails
		m.sendEmailDirect(queuedEmail)
	} else {
		logger.Logger.Info(fmt.Sprintf("Email to %s queued successfully", emailModel.To))
	}

	// Save notification record
	notification := entity.Notification{
		Type:      entity.EMAIL,
		Recipient: emailModel.To,
		Subject:   emailModel.Subject,
		Content:   emailModel.Message,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	_, err = m.notificationRepository.Create(ctx, notification)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to save email notification: %s", err.Error()))
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
		logger.Logger.Error(fmt.Sprintf("Error sending email directly: %s", err.Error()))
	} else {
		logger.Logger.Info("Email sent directly successfully!")
	}
}
