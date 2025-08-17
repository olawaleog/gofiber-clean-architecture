package model

type RatingModel struct {
	Review  string `json:"review"`
	Rating  uint   `json:"rating"`
	UserId  uint   `json:"userId"`
	OrderId uint   `json:"orderId"`
}
