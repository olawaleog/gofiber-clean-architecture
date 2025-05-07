package entity

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TransactionId   uint        `gorm:"column:transaction_id;type:int"`
	Transaction     Transaction `gorm:"foreignKey:TransactionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount          float64     `gorm:"column:amount;type:numeric(10,2)"`
	Currency        string      `gorm:"column:currency;type:varchar(50)"`
	WaterCost       float64     `gorm:"column:water_cost;type:numeric(10,2)"`
	DeliveryFee     float64     `gorm:"column:delivery_fee;type:numeric(10,2)"`
	DeliveryAddress string      `gorm:"column:delivery_address;type:varchar(255)"`
	DeliveryPlaceId string      `gorm:"column:delivery_place_id;type:varchar(255)"`
	RefineryAddress string      `gorm:"column:refinery_address;type:varchar(255)"`
	RefineryPlaceId string      `gorm:"column:refinery_place_id;type:varchar(255)"`
	RefineryId      uint        `gorm:"column:refinery_id;type:int"`
	Refinery        Refinery    `gorm:"foreignKey:RefineryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status          uint        `gorm:"column:status;type:int"`
	TruckId         uint        `gorm:"column:truck_id;type:int"`
	Capacity        string      `gorm:"column:capacity;type:varchar(50)"`
}

func (Order) TableName() string {
	return "tb_orders"
}
