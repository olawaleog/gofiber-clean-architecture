package model

type OtpModel struct {
	Code      string `json:"otp" validate:"required"`
	UserId    int    `json:"userId" validate:"required"`
	Operation string `json:"operationType" validate:"required"`
}
