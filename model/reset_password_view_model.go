package model

type ResetPasswordViewModel struct {
	UserId   int    `json:"userId"`
	OTP      string `json:"otp"`
	Password string `json:"password"`
}
