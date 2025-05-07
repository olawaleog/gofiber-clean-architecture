package model

import "time"

type OrderModel struct {
	TransactionId   uint             `json:"transactionId"`
	Transaction     TransactionModel `json:"transaction"`
	Amount          float64          `json:"amount"`
	Currency        string           `json:"currency"`
	WaterCost       float64          `json:"waterCost"`
	DeliveryFee     float64          `json:"deliveryFee"`
	DeliveryAddress string           `json:"deliveryAddress"`
	RefineryAddress string           `json:"refineryAddress"`
	RefineryId      uint             `json:"refineryId"`
	Status          uint             `json:"status"`
	TruckId         uint             `json:"truckId"`
	Capacity        string           `json:"capacity"`
	Id              uint             `json:"id"`
	CreatedAt       time.Time        `json:"createdAt"`
	Refinery        RefineryModel    `json:"refinery"`
	User            UserModel        `json:"user"`
	DeliveryPlaceId string           `json:"deliveryPlaceId"`
	RefineryPlaceId string           `json:"refineryPlaceId"`
}
