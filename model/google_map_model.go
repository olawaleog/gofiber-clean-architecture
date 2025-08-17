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
