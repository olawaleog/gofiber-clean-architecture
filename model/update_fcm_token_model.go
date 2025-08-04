package model

type UpdateFcmToken struct {
	Id       int64  `json:"id"`
	FcmToken string `json:"fcmToken"`
}
