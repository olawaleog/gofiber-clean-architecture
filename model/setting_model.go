package model

// SettingModel represents a setting in the system
type SettingModel struct {
	Key   string `json:"key" validate:"required,max=50"`
	Value string `json:"value" validate:"required"`
}

// SettingResponseModel is used for responses
type SettingResponseModel struct {
	ID    uint   `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
