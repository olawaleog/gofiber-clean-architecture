package model

type PlaceSuggestion struct {
	Description string `json:"description"`
	PlaceID     string `json:"place_id"`
}

type GeocodeResult struct {
	PlaceID   string `json:"place_id"`
	Address   string `json:"address"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// PlaceDetailResponse represents the response from Google Places API Details endpoint
type PlaceDetailResponse struct {
	Result struct {
		PlaceID          string `json:"place_id"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
		Name      string   `json:"name"`
		Icon      string   `json:"icon"`
		Types     []string `json:"types"`
		Vicinity  string   `json:"vicinity"`
		Reference string   `json:"reference"`
	} `json:"result"`
	Status string `json:"status"`
}
