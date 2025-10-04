package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type SettingController struct {
	SettingService service.SettingService
}

func NewSettingController(service service.SettingService) *SettingController {
	return &SettingController{SettingService: service}
}

func (controller *SettingController) Route(app *fiber.App) {
	app.Get("/v1/api/settings", controller.List)
	app.Get("/v1/api/settings/:key", controller.FindByKey)
	app.Post("/v1/api/settings", controller.Create)
	app.Put("/v1/api/settings/:key", controller.Update)
}

func (controller *SettingController) List(c *fiber.Ctx) error {
	settings, err := controller.SettingService.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "Error fetching settings",
			Data:    err.Error(),
		})
	}

	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    settings,
	})
}

func (controller *SettingController) FindByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	setting, err := controller.SettingService.FindByKey(c.Context(), key)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.GeneralResponse{
			Code:    404,
			Message: "Setting not found",
			Data:    err.Error(),
		})
	}

	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    setting,
	})
}

func (controller *SettingController) Create(c *fiber.Ctx) error {
	var request model.SettingModel
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Bad request",
			Data:    err.Error(),
		})
	}

	setting, err := controller.SettingService.Create(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "Error creating setting",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    201,
		Message: "Setting created successfully",
		Data: model.SettingResponseModel{
			ID:    setting.ID,
			Key:   setting.Key,
			Value: setting.Value,
		},
	})
}

func (controller *SettingController) Update(c *fiber.Ctx) error {
	key := c.Params("key")
	var request model.SettingModel
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: "Bad request",
			Data:    err.Error(),
		})
	}

	setting, err := controller.SettingService.Update(c.Context(), key, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    500,
			Message: "Error updating setting",
			Data:    err.Error(),
		})
	}

	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Setting updated successfully",
		Data: model.SettingResponseModel{
			ID:    setting.ID,
			Key:   setting.Key,
			Value: setting.Value,
		},
	})
}
