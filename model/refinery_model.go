package model

import "time"

type RefineryModel struct {
	Id                             uint      `json:"id"`
	Name                           string    `json:"name"`
	Address                        string    `json:"address"`
	Phone                          string    `json:"phone"`
	Email                          string    `json:"email"`
	Website                        string    `json:"website"`
	Licence                        string    `json:"licence"`
	LicenceExpiry                  time.Time `json:"licenceExpiry"`
	PlaceId                        string    `json:"placeId"`
	RawLocationData                string    `json:"rawLocationData"`
	Region                         string    `json:"region"`
	DomesticCostPerThousandLitre   float64   `json:"domesticCostPerThousandLitre"`
	IndustrialCostPerThousandLitre float64   `json:"industrialCostPerThousandLitre"`
	IsActive                       bool      `json:"isActive"`
	Longitude                      string    `json:"longitude"`
	Latitude                       string    `json:"latitude"`
	HasDomesticWaterSupply         bool      `json:"hasDomesticWaterSupply"`
	HasIndustrialWaterSupply       bool      `json:"hasIndustrialWaterSupply"`
}
