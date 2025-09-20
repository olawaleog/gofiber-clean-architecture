package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type PaymentConfigController struct {
	service.PaymentConfigService
}

func NewPaymentConfigController(service service.PaymentConfigService) *PaymentConfigController {
	return &PaymentConfigController{PaymentConfigService: service}
}

func (controller *PaymentConfigController) Route(app *fiber.App) {
	app.Get("/v1/api/payment-configs", controller.List)
	app.Get("/v1/v1/api/payment-configs/:countryCode", controller.Get)
	app.Post("/v1/api/payment-configs", controller.Create)
	app.Put("/v1/api/payment-configs/:countryCode", controller.Update)
	app.Delete("/v1/api/payment-configs/:countryCode", controller.Delete)
}

func (controller *PaymentConfigController) List(c *fiber.Ctx) error {
	configs, err := controller.PaymentConfigService.ListPaymentConfigs(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "Error fetching payment configurations",
			Data:    err.Error(),
		})
	}

	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    configs,
	})
}

func (controller *PaymentConfigController) Get(c *fiber.Ctx) error {
	countryCode := c.Params("countryCode")
	config, err := controller.PaymentConfigService.GetPaymentConfig(c.Context(), countryCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.GeneralResponse{
			Code:    404,
			Message: "Payment configuration not found",
			Data:    err.Error(),
		})
	}

	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    config,
	})
}

func (controller *PaymentConfigController) Create(c *fiber.Ctx) error {
	var config model.PaymentConfigModel
	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	if err := controller.PaymentConfigService.CreatePaymentConfig(c.Context(), &config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "Error creating payment configuration",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    201,
		Message: "Payment configuration created successfully",
		Data:    config,
	})
}

func (controller *PaymentConfigController) Update(c *fiber.Ctx) error {
	countryCode := c.Params("countryCode")
	var config model.PaymentConfigModel
	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	config.CountryCode = countryCode
	if err := controller.PaymentConfigService.UpdatePaymentConfig(c.Context(), &config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "Error updating payment configuration",
			Data:    err.Error(),
		})
	}

	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Payment configuration updated successfully",
		Data:    config,
	})
}

func (controller *PaymentConfigController) Delete(c *fiber.Ctx) error {
	countryCode := c.Params("countryCode")
	if err := controller.PaymentConfigService.DeletePaymentConfig(c.Context(), countryCode); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "Error deleting payment configuration",
			Data:    err.Error(),
		})
	}

	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Payment configuration deleted successfully",
		Data:    nil,
	})
}
