package model

type MobileMoneyRequestModel struct {
	PaymentMethodId uint    `json:"paymentMethodId"`
	Amount          float64 `json:"amount"`
	WaterCost       float64 `json:"waterCost"`
	DeliveryFee     float64 `json:"deliveryFee"`
	Provider        string  `json:"provider"`
	EmailAddress    string  `json:"emailAddress"`
	PhoneNumber     string  `json:"phoneNumber"`
	UserId          uint    `json:"userId"`
	Currency        string  `json:"currency"`
	CustomerPlaceId string  `json:"customerPlaceId"`
	RefineryPlaceId string  `json:"refineryPlaceId"`
	CustomerAddress string  `json:"customerAddress"`
	RefineryAddress string  `json:"refineryAddress"`
	RefineryId      uint    `json:"refineryId"`
	Capacity        string  `json:"capacity"`
	Type            string  `json:"type"`
	AddressId       uint    `json:"addressId"`
	CountryCode     string  `json:"countryCode"`
}
