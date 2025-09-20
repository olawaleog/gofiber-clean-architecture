package model

type PaymentConfigModel struct {
	ID          uint   `json:"id"`
	CountryCode string `json:"countryCode"`
	SecretKey   string `json:"secretKey"`
	PublicKey   string `json:"publicKey"`
}
