package impl

import (
	"encoding/json"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"gopkg.in/gomail.v2"
	"strconv"
)

// EmailConsumerServiceImpl implements the EmailConsumerService interface
type EmailConsumerServiceImpl struct {
	config          configuration.Config
	rabbitMQService *RabbitMQService
}

// NewEmailConsumerService creates a new EmailConsumerService
func NewEmailConsumerService(config configuration.Config, rabbitMQService *RabbitMQService) service.EmailConsumerService {
	return &EmailConsumerServiceImpl{
		config:          config,
		rabbitMQService: rabbitMQService,
	}
}

// StartConsumer starts consuming email messages from the queue
func (e *EmailConsumerServiceImpl) StartConsumer() error {
	common.Logger.Info("Starting email consumer service...")

	// Subscribe to the email.send routing key
	err := e.rabbitMQService.SubscribeToTopic("email.send", e.ProcessEmail)
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Failed to start email consumer: %s", err.Error()))
		return err
	}

	common.Logger.Info("Email consumer service started successfully")
	return nil
}

// ProcessEmail processes a single email message
func (e *EmailConsumerServiceImpl) ProcessEmail(emailData []byte) error {
	// Parse the email data
	var queuedEmail model.QueuedEmailMessage
	err := json.Unmarshal(emailData, &queuedEmail)
	if err != nil {
		common.Logger.Error(fmt.Sprintf("Error parsing email data: %s", err.Error()))
		return err
	}

	common.Logger.Info(fmt.Sprintf("Processing email to %s with subject: %s", queuedEmail.To, queuedEmail.Subject))

	// Create a new email message
	message := gomail.NewMessage()

	// Set email headers
	from := queuedEmail.From
	if from == "" {
		from = e.config.Get("MAIL_FROM")
	}
	message.SetHeader("From", from)
	message.SetHeader("To", queuedEmail.To)
	message.SetHeader("Subject", queuedEmail.Subject)

	// Set email body
	message.SetBody("text/html", queuedEmail.Message)

	// Set up the SMTP connection
	port, _ := strconv.Atoi(e.config.Get("MAIL_PORT"))
	dialer := gomail.NewDialer(
		e.config.Get("MAIL_HOST"),
		port,
		e.config.Get("MAIL_USERNAME"),
		e.config.Get("MAIL_PASSWORD"),
	)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		common.Logger.Error(fmt.Sprintf("Error sending email: %s", err.Error()))
		return err
	}

	common.Logger.Info(fmt.Sprintf("Email to %s sent successfully", queuedEmail.To))
	return nil
}
