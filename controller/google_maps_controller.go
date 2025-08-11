package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type GoogleMapsController struct {
	mapsService service.GoogleMapsService
}

func NewGoogleMapsController(mapsService service.GoogleMapsService) *GoogleMapsController {
	return &GoogleMapsController{mapsService: mapsService}
}

func (g *GoogleMapsController) RegisterRoutes(router fiber.Router) {
	router.Get("/v1/api/googlemaps/suggest", g.PlaceSuggestion)
	router.Get("/v1/api/googlemaps/reverse-geocode", g.ReverseGeocode)
}

func (g *GoogleMapsController) PlaceSuggestion(c *fiber.Ctx) error {
	input := c.Query("input")
	if input == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input query param required"})
	}
	suggestions, err := g.mapsService.SuggestPlaces(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(suggestions)
}

func (g *GoogleMapsController) ReverseGeocode(c *fiber.Ctx) error {
	lat := c.Query("lat")
	lng := c.Query("lng")
	if lat == "" || lng == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "lat and lng query params required"})
	}
	result, err := g.mapsService.ReverseGeocode(c.Context(), lat, lng)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}
