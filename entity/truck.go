package entity

import (
	"gorm.io/gorm"
	"time"
)

type Truck struct {
	gorm.Model
	ManufacturerModel     string    `gorm:"column:manufacturer_model;type:varchar(100)"`
	PlateNumber           string    `gorm:"column:plate_number;type:varchar(100)"`
	TruckType             string    `gorm:"column:truck_type;type:varchar(100)"`
	IsActive              bool      `gorm:"column:is_active;type:boolean"`
	Capacity              int       `gorm:"column:capacity;type:int"`
	YearOfManufacture     int       `gorm:"column:year_of_manufacture;type:int"`
	UserId                uint      `gorm:"column:user_id;type:int"`
	EngineNumber          string    `gorm:"column:engine_number;type:varchar(100)"`
	LicenceExpirationDate time.Time `gorm:"column:licence_expiration_date;type:date"`
	User                  User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Truck) TableName() string {
	return "tb_trucks"
}
