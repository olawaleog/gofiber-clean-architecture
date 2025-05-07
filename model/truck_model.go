package model

import "time"

type TruckModel struct {
	Id                    uint      `json:"id"`
	YearOfManufacture     string    `json:"yearOfManufacture"`
	Capacity              string    `json:"capacity"`
	IsActive              bool      `json:"isActive"`
	ManufacturerModel     string    `json:"model"`
	PlateNumber           string    `json:"plateNumber"`
	EngineNumber          string    `json:"engineNumber"`
	LicenceExpirationDate time.Time `json:"licenceExpirationDate"`

	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Phone     string    `json:"phoneNumber"`
	Email     string    `json:"emailAddress"`
	UserId    uint      `json:"userId"`
	User      UserModel `json:"user"`
}
