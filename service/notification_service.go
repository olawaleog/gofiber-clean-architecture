package service

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type NotificationService interface {
	// SendToDevice sends a notification to a single device
	SendToDevice(ctx context.Context, notification model.NotificationModel) error

	// SendToMultipleDevices sends a notification to multiple devices
	SendToMultipleDevices(ctx context.Context, notification model.NotificationModel) error

	// SendToTopic sends a notification to a topic
	SendToTopic(ctx context.Context, topic string, notification model.NotificationModel) error

	// SubscribeToTopic subscribes devices to a topic
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) error

	// UnsubscribeFromTopic unsubscribes devices from a topic
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) error
}
