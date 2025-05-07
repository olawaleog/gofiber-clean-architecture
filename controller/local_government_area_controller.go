package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type LocalGovernmentAreaController struct {
	service.LocalGovernmentService
}

func NewLocalGovernmentAreaController(service *service.LocalGovernmentService) LocalGovernmentAreaController {
	return LocalGovernmentAreaController{LocalGovernmentService: *service}
}

func (l LocalGovernmentAreaController) Route(app fiber.Router) {
	app.Get("/v1/api/local-government-areas", l.FindAll)
	app.Put("/v1/api/local-government-areas/:id/toggle-active", l.ToggleLocalGovernmentActive)
	app.Get("/v1/api/get-place-suggestion/:id", l.GetPlaceSuggestion)
}

func (l LocalGovernmentAreaController) FindAll(c *fiber.Ctx) error {
	localGovernmentAres, err := l.LocalGovernmentService.FindAll(c.Context())
	exception.PanicLogging(err)
	return c.Status(fiber.StatusOK).JSON(localGovernmentAres)
}

func (l LocalGovernmentAreaController) ToggleLocalGovernmentActive(c *fiber.Ctx) error {
	id := c.Params("id")
	err := l.LocalGovernmentService.ToggleLocalGovernmentActive(c.Context(), id)
	exception.PanicLogging(err)
	return c.Status(fiber.StatusOK).JSON(nil)
}

func (l LocalGovernmentAreaController) GetPlaceSuggestion(c *fiber.Ctx) error {
	id := c.Params("id")
	response := l.LocalGovernmentService.GetPlaceSuggestion(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(
		model.GeneralResponse{
			Code:    200,
			Message: "Success",
			Data:    response,
			Success: true,
		})
}
