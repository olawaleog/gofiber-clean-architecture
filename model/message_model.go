package model

type EmailMessageModel struct {
	Subject  string
	To       string
	Message  string
	File     []byte
	FileName string
}

type SMSMessageModel struct {
	PhoneNumber string `json:"phone_number" json:"phoneNumber"`
	Message     string `json:"message"`
	CountryCode string `json:"country_code" json:"countryCode"`
}
