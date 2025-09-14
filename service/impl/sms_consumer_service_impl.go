package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
)

type SMSConsumerService interface {
	StartConsumer() error
}

type SMSConsumerServiceImpl struct {
	httpService            service.HttpService
	rabbitMQService        *RabbitMQService
	config                 map[string]string
	notificationRepository repository.NotificationRepository
}

func NewSMSConsumerService(httpService service.HttpService, rabbitMQService *RabbitMQService,
	config map[string]string, notificationRepository repository.NotificationRepository) SMSConsumerService {
	return &SMSConsumerServiceImpl{
		httpService:            httpService,
		rabbitMQService:        rabbitMQService,
		config:                 config,
		notificationRepository: notificationRepository,
	}
}

func (s *SMSConsumerServiceImpl) StartConsumer() error {
	logger.Logger.Info("Starting SMS consumer service")
	return s.rabbitMQService.SubscribeToTopic("sms.send", func(message []byte) error {
		logger.Logger.Info("Received SMS message to send")

		var smsMessage model.QueuedSMSMessage
		if err := json.Unmarshal(message, &smsMessage); err != nil {
			logger.Logger.Error(fmt.Sprintf("Failed to unmarshal SMS message: %s", err.Error()))
			return err
		}

		// Process the SMS message
		err := s.processSMSMessage(smsMessage)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("Failed to process SMS message: %s", err.Error()))
			return err
		}

		return nil
	})
}

func (s *SMSConsumerServiceImpl) processSMSMessage(smsMessage model.QueuedSMSMessage) error {
	headers := make(map[string]interface{})
	headers["apiKey"] = s.config["AFRICAS_TALKING_API_KEY"]

	body := make(map[string]interface{})
	body["username"] = s.config["AFRICAS_TALKING_USERNAME"]
	body["enqueue"] = 1
	body["message"] = smsMessage.Message

	phoneNumber := smsMessage.PhoneNumber
	if phoneNumber[0] == '0' {
		phoneNumber = smsMessage.CountryCode + phoneNumber[1:]
	} else if smsMessage.CountryCode != "" && phoneNumber[0] != '+' {
		phoneNumber = smsMessage.CountryCode + phoneNumber
	}

	body["to"] = phoneNumber
	body["from"] = "AquaWizz"

	url := s.config["AFRICAS_TALKING_BASE_URL"] + "/version1/messaging"
	_, err := s.httpService.PostMethod(context.Background(), url, "POST", &body, &headers, true)

	// Update the notification status based on the result
	status := "delivered"
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error sending SMS via Africa's Talking API: %s", err.Error()))
		status = "failed"
	} else {
		logger.Logger.Info(fmt.Sprintf("SMS sent to %s successfully", phoneNumber))
	}

	// Update notification record if user ID is provided
	if smsMessage.UserID > 0 {
		notification := entity.Notification{
			UserID:    smsMessage.UserID,
			Type:      entity.SMS,
			Recipient: phoneNumber,
			Content:   smsMessage.Message,
			Status:    status,
			SentAt:    time.Now(),
		}

		_, saveErr := s.notificationRepository.Create(context.Background(), notification)
		if saveErr != nil {
			logger.Logger.Error(fmt.Sprintf("Failed to save SMS notification: %s", saveErr.Error()))
			// We don't return this error as the SMS was already sent or attempted,
			// this is just for record keeping
		}
	}

	return nil
}
