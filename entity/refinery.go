package entity

import (
	"time"

	"gorm.io/gorm"
)

type Refinery struct {
	gorm.Model
	Name                           string    `gorm:";column:name;type:text"`
	Address                        string    `gorm:"column:address;type:text"`
	Licence                        string    `gorm:"column:licence;type:text"`
	Phone                          string    `gorm:"column:phone;type:text"`
	Email                          string    `gorm:"column:email;type:text"`
	Website                        string    `gorm:"column:website;type:text"`
	LicenceExpiry                  time.Time `gorm:"column:licence_expiry;type:timestamp"`
	PlaceId                        string    `gorm:"column:place_id;type:text"`
	RawLocationData                string    `gorm:"column:raw_location_data;type:text"`
	Region                         string    `gorm:"column:region;type:text"`
	Country                        string    `gorm:"column:country;type:text"`
	Currency                       string    `gorm:"column:currency;type:text"`
	CurrencySymbol                 string    `gorm:"column:currency_symbol;type:text"`
	DomesticCostPerThousandLitre   float64   `gorm:"column:domestic_cost_per_thousand_litre;type:numeric(10,2)"`
	IndustrialCostPerThousandLitre float64   `gorm:"column:industrial_cost_per_thousand_litre;type:numeric(10,2)"`
	Longitude                      string    `gorm:"column:longitude;type:text"`
	Latitude                       string    `gorm:"column:latitude;type:text"`
	HasDomesticWaterSupply         bool      `gorm:"column:has_domestic_water_supply;type:boolean"`
	HasIndustrialWaterSupply       bool      `gorm:"column:has_industrial_water_supply;type:boolean"`
	IsActive                       bool      `gorm:"column:is_active;type:boolean"`
}

func (Refinery) TableName() string {
	return "tb_refineries"
}
