package model

type PaymentMethodModel struct {
	UniqueId string `json:"uniqueId"`
	Provider string `json:"provider"`
	Scheme   string `json:"scheme"`
	AuthCode string `json:"authCode"`
	Name     string `json:"name"`
	Id       uint   `json:"id"`
}
