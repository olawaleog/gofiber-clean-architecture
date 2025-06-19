package impl

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"google.golang.org/api/option"
)

type notificationServiceImpl struct {
	app *firebase.App
}

func NewNotificationService(credentialsPath string) service.NotificationService {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	exception.PanicLogging(err)

	return &notificationServiceImpl{
		app: app,
	}
}

func (n *notificationServiceImpl) SendToDevice(ctx context.Context, notification model.NotificationModel) error {
	if notification.Token == "" {
		return fmt.Errorf("device token is required")
	}

	client, err := n.app.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.ImageURL,
		},
		Token: notification.Token,
		Data:  notification.Data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ClickAction: notification.ClickAction,
			},
		},
	}

	if notification.ClickAction != "" {
		message.APNS = &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: notification.Title,
						Body:  notification.Body,
					},
				},
			},
		}
	}

	_, err = client.Send(ctx, message)
	return err
}

func (n *notificationServiceImpl) SendToMultipleDevices(ctx context.Context, notification model.NotificationModel) error {
	if len(notification.Tokens) == 0 {
		return fmt.Errorf("at least one device token is required")
	}

	client, err := n.app.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.ImageURL,
		},
		Tokens: notification.Tokens,
		Data:   notification.Data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ClickAction: notification.ClickAction,
			},
		},
	}

	if notification.ClickAction != "" {
		message.APNS = &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: notification.Title,
						Body:  notification.Body,
					},
				},
			},
		}
	}

	_, err = client.SendMulticast(ctx, message)
	return err
}

func (n *notificationServiceImpl) SendToTopic(ctx context.Context, topic string, notification model.NotificationModel) error {
	if topic == "" {
		return fmt.Errorf("topic is required")
	}

	client, err := n.app.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.ImageURL,
		},
		Topic: topic,
		Data:  notification.Data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ClickAction: notification.ClickAction,
			},
		},
	}

	if notification.ClickAction != "" {
		message.APNS = &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: notification.Title,
						Body:  notification.Body,
					},
				},
			},
		}
	}

	_, err = client.Send(ctx, message)
	return err
}

func (n *notificationServiceImpl) SubscribeToTopic(ctx context.Context, tokens []string, topic string) error {
	if len(tokens) == 0 {
		return fmt.Errorf("at least one device token is required")
	}
	if topic == "" {
		return fmt.Errorf("topic is required")
	}

	client, err := n.app.Messaging(ctx)
	if err != nil {
		return err
	}

	_, err = client.SubscribeToTopic(ctx, tokens, topic)
	return err
}

func (n *notificationServiceImpl) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) error {
	if len(tokens) == 0 {
		return fmt.Errorf("at least one device token is required")
	}
	if topic == "" {
		return fmt.Errorf("topic is required")
	}

	client, err := n.app.Messaging(ctx)
	if err != nil {
		return err
	}

	_, err = client.UnsubscribeFromTopic(ctx, tokens, topic)
	return err
}
