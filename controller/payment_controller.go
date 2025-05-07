package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	service.PaymentService
	service.UserService
	configuration.Config
}

func (c PaymentController) Route(app *fiber.App) {
	app.Get("/v1/api/get-payment-options", c.GetPaymentMethods)

}

func (c PaymentController) GetPaymentMethods(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	paymentMethods := c.PaymentService.GetPaymentMethods(ctx.Context(), claims["userId"].(float64))
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    paymentMethods,
		Success: true,
	})
}

func NewPaymentController(paymentService *service.PaymentService, userService *service.UserService, config configuration.Config) *PaymentController {
	return &PaymentController{
		PaymentService: *paymentService,
		UserService:    *userService,
		Config:         config,
	}
}
