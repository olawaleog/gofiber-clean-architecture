package model

type PlaceDetailResponse struct {
	HtmlAttributions []interface{} `json:"html_attributions"`
	Result           struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"result"`
	Status string
}
