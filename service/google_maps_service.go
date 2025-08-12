package service

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type GoogleMapsService interface {
	SuggestPlaces(ctx context.Context, input string) ([]model.PlaceSuggestion, error)
	ReverseGeocode(ctx context.Context, lat, lng string) (*model.GeocodeResult, error)
	GetPlaceDetail(ctx context.Context, id string) (*model.PlaceDetailResponse, error)
}

type googleMapsService struct {
	apiKey string
}

//type PlaceSuggestion struct {
//	Description string `json:"description"`
//	PlaceID     string `json:"place_id"`
//}
//
//type GeocodeResult struct {
//	PlaceID string `json:"place_id"`
//	Address string `json:"address"`
//}
//
//func NewGoogleMapsService() GoogleMapsService {
//	return &googleMapsService{
//		apiKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
//	}
//}
//
//func (g *googleMapsService) SuggestPlaces(ctx context.Context, input string) ([]PlaceSuggestion, error) {
//	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/autocomplete/json?input=%s&key=%s", input, g.apiKey)
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	var data struct {
//		Predictions []struct {
//			Description string `json:"description"`
//			PlaceID     string `json:"place_id"`
//		} `json:"predictions"`
//	}
//	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
//		return nil, err
//	}
//
//	suggestions := make([]PlaceSuggestion, len(data.Predictions))
//	for i, p := range data.Predictions {
//		suggestions[i] = PlaceSuggestion{Description: p.Description, PlaceID: p.PlaceID}
//	}
//	return suggestions, nil
//}
//
//func (g *googleMapsService) ReverseGeocode(ctx context.Context, lat, lng string) (*GeocodeResult, error) {
//	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%s,%s&key=%s", lat, lng, g.apiKey)
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	var data struct {
//		Results []struct {
//			PlaceID          string `json:"place_id"`
//			FormattedAddress string `json:"formatted_address"`
//		} `json:"results"`
//	}
//	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
//		return nil, err
//	}
//	if len(data.Results) == 0 {
//		return nil, fmt.Errorf("no results found")
//	}
//	return &GeocodeResult{
//		PlaceID: data.Results[0].PlaceID,
//		Address: data.Results[0].FormattedAddress,
//	}, nil
//}
