package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
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
	router.Get("/v1/api/googlemaps/place-detail", g.PlaceDetail)
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
	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    suggestions,
		Success: true,
	})
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
	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
		Success: true,
	})
}

func (g *GoogleMapsController) PlaceDetail(c *fiber.Ctx) error {
	placeID := c.Query("place_id")
	if placeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "place_id query param required"})
	}
	detail, err := g.mapsService.GetPlaceDetail(c.Context(), placeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    detail,
		Success: true,
	})
}
