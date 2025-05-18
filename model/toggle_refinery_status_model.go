package model

type ToggleRefineryStatusModel struct {
	Id     uint `json:"id" validate:"required"`
	Status bool `json:"status" validate:"required"`
}
