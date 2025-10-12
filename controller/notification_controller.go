package controller

import (
	"context"
	"net/http"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"      // Added import for NotificationModel
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository" // Added import for UserRepository
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type NotificationController struct {
	NotificationService service.NotificationService
	UserRepository      repository.UserRepository // Fixed incorrect type
}

func NewNotificationController(notificationService service.NotificationService, userRepository repository.UserRepository) *NotificationController {
	return &NotificationController{
		NotificationService: notificationService,
		UserRepository:      userRepository,
	}
}

func (c *NotificationController) SendWelcomeNotification(ctx *fiber.Ctx) error {
	// Fetch all users with FCM tokens
	users, err := c.UserRepository.FindAllWithFcmToken(context.Background())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	// Extract FCM tokens
	tokens := []string{}
	for _, user := range users {
		if user.FcmToken != "" {
			tokens = append(tokens, user.FcmToken)
		}
	}

	// Send notification
	err = c.NotificationService.SendToMultipleDevices(context.Background(), model.NotificationModel{
		Tokens: tokens,
		Title:  "Welcome to Aqua Wizz",
		Body:   "Welcome to Aqua Wizz",
	})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send notifications",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Notifications sent successfully",
	})
}

func (c *NotificationController) Route(app *fiber.App) {
	app.Post("/notifications/welcome", c.SendWelcomeNotification)
}
