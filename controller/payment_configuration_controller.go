package controller

import (
	"context"
	"net/http"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type PaymentConfigurationController struct {
	PaymentConfigService service.PaymentConfigService
}

func NewPaymentConfigurationController(paymentConfigService service.PaymentConfigService) *PaymentConfigurationController {
	return &PaymentConfigurationController{
		PaymentConfigService: paymentConfigService,
	}
}

func (c *PaymentConfigurationController) GetPaymentConfiguration(ctx *fiber.Ctx) error {
	countryCode := ctx.Query("countryCode")
	config, err := c.PaymentConfigService.GetPaymentConfig(context.Background(), countryCode)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(config)
}

func (c *PaymentConfigurationController) RefreshCache(ctx *fiber.Ctx) error {
	if err := c.PaymentConfigService.RefreshCache(context.Background()); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.SendStatus(http.StatusOK)
}
