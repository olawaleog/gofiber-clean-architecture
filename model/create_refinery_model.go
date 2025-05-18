package model

import "time"

type CreateRefineryModel struct {
	Name                           string    `json:"name"`
	PlaceId                        string    `json:"placeId"`
	Address                        string    `json:"address"`
	Longitude                      string    `json:"longitude"`
	Latitude                       string    `json:"latitude"`
	HasDomesticWaterSupply         bool      `json:"hasDomesticWaterSupply"`
	HasIndustrialWaterSupply       bool      `json:"hasIndustrialWaterSupply"`
	LicenceExpirationDate          time.Time `json:"licenceExpirationDate"`
	Capacity                       string    `json:"capacity"`
	DomesticCostPerThousandLitre   float64   `json:"domesticCostPerThousandLitre"`
	IndustrialCostPerThousandLitre float64   `json:"industrialCostPerThousandLitre"`
	RawLocationData                string    `json:"rawLocationData"`
	Region                         string    `json:"region"`
	Phone                          string    `json:"phone"`
	Email                          string    `json:"email"`
	Website                        string    `json:"website"`
	FirstName                      string    `json:"firstName"`
	LastName                       string    `json:"lastName"`
}
