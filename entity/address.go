package entity

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type Address struct {
	gorm.Model
	UserId       uint   `gorm:"column:user_id"`
	Street       string `gorm:"column:street;type:varchar(100)"`
	City         string `gorm:"column:city;type:varchar(100)"`
	PostalCode   string `gorm:"column:postal_code;type:varchar(10)"`
	IsMain       bool   `gorm:"column:is_main;type:boolean"`
	Longitude    string `gorm:"column:longitude;type:varchar(100)"`
	Latitude     string `gorm:"column:latitude;type:varchar(100)"`
	Region       string `gorm:"column:region;type:varchar(100)"`
	Raw          JSONB  `gorm:"column:raw;type:jsonb"`
	Description  string `gorm:"column:description;type:varchar(100)"`
	StreetNumber string `gorm:"column:street_number;type:varchar(100)"`
	PlaceId      string `gorm:"column:place_id;type:varchar(100)"`
	CountryCode  string `gorm:"column:country_code;type:varchar(10)"`
}

func (Address) TableName() string {
	return "tb_addresses"
}
