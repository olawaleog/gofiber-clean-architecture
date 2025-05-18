package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type RefineryController struct {
	service.RefineryService
	service.UserService
	configuration.Config
}

func (controller RefineryController) Route(app *fiber.App) {
	app.Post("/v1/api/refinery/get", controller.GetRefinery)
	app.Post("/v1/api/refinery/create", controller.CreateRefinery)
	app.Post("/v1/api/refinery/update/:id", controller.UpdateRefinery)
	app.Get("/v1/api/refinery/list", controller.ListRefineries)
	app.Post("/v1/api/refinery/toggle-status", controller.ToggleRefineryStatus)

}

func NewRefineryController(refineryService *service.RefineryService, config configuration.Config) *RefineryController {
	return &RefineryController{
		Config:          config,
		RefineryService: *refineryService,
	}
}

func (controller RefineryController) GetRefinery(ctx *fiber.Ctx) error {
	var getRefineryModel model.GetRefineryModel
	err := ctx.BodyParser(&getRefineryModel)
	exception.PanicLogging(err)
	refinery, err := controller.RefineryService.GetRefinery(ctx.Context(), getRefineryModel)

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    refinery,
	})
}

func (controller RefineryController) CreateRefinery(ctx *fiber.Ctx) error {
	var createRefineryModel model.CreateRefineryModel

	err := ctx.BodyParser(&createRefineryModel)
	exception.PanicLogging(err)

	refinery, err := controller.RefineryService.CreateRefinery(ctx.Context(), createRefineryModel)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    fiber.StatusCreated,
		Message: "Successful",
		Data:    refinery,
	})
}

func (controller RefineryController) ListRefineries(ctx *fiber.Ctx) error {
	refinaries, err := controller.RefineryService.ListRefineries(ctx.Context())
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    refinaries,
		Success: true,
	})
}

func (controller RefineryController) UpdateRefinery(ctx *fiber.Ctx) error {
	var updateRefineryModel model.CreateRefineryModel
	err := ctx.BodyParser(&updateRefineryModel)
	exception.PanicLogging(err)

	id := ctx.Params("id")

	refinery, err := controller.RefineryService.UpdateRefinery(ctx.Context(), updateRefineryModel, id)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    refinery,
	})
}

func (controller RefineryController) ToggleRefineryStatus(ctx *fiber.Ctx) error {
	var toggleRefineryStatusModel model.ToggleRefineryStatusModel
	err := ctx.BodyParser(&toggleRefineryStatusModel)
	exception.PanicLogging(err)

	refinery, err := controller.RefineryService.ToggleRefineryStatus(ctx.Context(), toggleRefineryStatusModel)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    refinery,
	})
}
