package impl

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
)

type GoogleMapsService interface {
	SuggestPlaces(ctx context.Context, input string) ([]model.PlaceSuggestion, error)
	ReverseGeocode(ctx context.Context, lat, lng string) (*model.GeocodeResult, error)
	GetPlaceDetail(ctx context.Context, placeID string) (*model.PlaceDetailResponse, error)
}

type googleMapsService struct {
	apiKey        string
	signingSecret string
}

func NewGoogleMapsService(configuration configuration.Config) service.GoogleMapsService {
	return &googleMapsService{
		apiKey:        configuration.Get("GOOGLE_MAPS_API_KEY"),
		signingSecret: configuration.Get("GOOGLE_MAPS_SIGNING_SECRET"),
	}
}

func (g *googleMapsService) signURL(rawURL string) (string, error) {
	if g.signingSecret == "" {
		return rawURL, nil // No signing secret set, skip signing
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	// Only sign the path and query
	pathAndQuery := u.EscapedPath()
	if u.RawQuery != "" {
		pathAndQuery += "?" + u.RawQuery
	}
	// Decode the secret
	secret, err := base64.URLEncoding.DecodeString(g.signingSecret)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha1.New, secret)
	mac.Write([]byte(pathAndQuery))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	// Add signature to query
	q := u.Query()
	q.Set("signature", signature)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (g *googleMapsService) SuggestPlaces(ctx context.Context, input string) ([]model.PlaceSuggestion, error) {
	urlStr := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/autocomplete/json?input=%s&key=%s&components=country:gh|country:ng&radius=1000000", input, g.apiKey)
	signedURL, err := g.signURL(urlStr)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(signedURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// To get the response body as a string:
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // or handle error as appropriate
	}
	bodyString := string(bodyBytes)
	// Now bodyString contains the full response body as a string
	fmt.Println(bodyString)
	var data struct {
		Predictions []struct {
			Description string `json:"description"`
			PlaceID     string `json:"place_id"`
		} `json:"predictions"`
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}

	suggestions := make([]model.PlaceSuggestion, len(data.Predictions))
	for i, p := range data.Predictions {
		suggestions[i] = model.PlaceSuggestion{Description: p.Description, PlaceID: p.PlaceID}
	}
	return suggestions, nil
}

func (g *googleMapsService) ReverseGeocode(ctx context.Context, lat, lng string) (*model.GeocodeResult, error) {
	urlStr := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%s,%s&key=%s", lat, lng, g.apiKey)

	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// To get the response body as a string:
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // or handle error as appropriate
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	// Now bodyString contains the full response body as a string

	var data struct {
		Results []struct {
			PlaceID          string `json:"place_id"`
			FormattedAddress string `json:"formatted_address"`
		} `json:"results"`
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}
	if len(data.Results) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	return &model.GeocodeResult{
		PlaceID:   data.Results[0].PlaceID,
		Address:   data.Results[0].FormattedAddress,
		Longitude: lng,
		Latitude:  lat,
	}, nil
}

func (g *googleMapsService) GetPlaceDetail(ctx context.Context, placeID string) (*model.PlaceDetailResponse, error) {
	urlStr := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&key=%s", placeID, g.apiKey)
	signedURL, err := g.signURL(urlStr)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(signedURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var detail model.PlaceDetailResponse
	if err := json.Unmarshal(bodyBytes, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}
