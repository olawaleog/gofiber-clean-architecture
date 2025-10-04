package model

type GetRefineryModel struct {
	AddressId   uint   `json:"addressId"`
	Type        string `json:"type"`
	CountryCode string `json:"countryCode" validate:"required,len=2"`
}
