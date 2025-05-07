package model

type WaterCostModel struct {
	WaterCost   float64 `json:"waterCost"`
	DeliveryFee float64 `json:"deliveryFee"`
	TotalCost   float64 `json:"totalCost"`
}
