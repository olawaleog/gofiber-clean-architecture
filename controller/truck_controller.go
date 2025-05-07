package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type TruckController struct {
	service.TruckService
	configuration.Config
}

func NewTruckController(s *service.TruckService, c configuration.Config) TruckController {
	return TruckController{TruckService: *s, Config: c}
}

func (controller TruckController) Route(app *fiber.App) {
	app.Get("/v1/api/trucks", controller.ListAllTrucks)
	app.Post("/v1/api/truck", controller.CreateTruck)
	app.Put("/v1/api/truck", controller.UpdateTruck)

}

func (controller TruckController) ListAllTrucks(c *fiber.Ctx) error {
	trucks, err := controller.TruckService.ListAllTrucks(c.Context())
	exception.PanicLogging(err)
	return c.Status(fiber.StatusOK).JSON(trucks)
}
func (controller TruckController) CreateTruck(c *fiber.Ctx) error {
	var request model.TruckModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	truck, err := controller.TruckService.CreateTruck(request)
	exception.PanicLogging(err)

	return c.Status(fiber.StatusOK).JSON(truck)
}

func (controller TruckController) UpdateTruck(c *fiber.Ctx) error {
	var request model.TruckModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	truck, err := controller.TruckService.UpdateTruck(request)
	exception.PanicLogging(err)

	return c.Status(fiber.StatusOK).JSON(truck)
}
