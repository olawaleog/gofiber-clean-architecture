package model

type ApproveOrRejectOrderModel struct {
	OrderId uint   `json:"orderId" validate:"required"`
	Action  string `json:"action" validate:"required"`
	TruckId uint   `json:"truckId""`
}
