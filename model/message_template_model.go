package model

type MessageTemplateModel struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	IsEmail bool   `json:"isEmail"`
	IsSms   bool   `json:"isSms"`
}
