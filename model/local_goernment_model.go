package model

type LocalGovernmentModel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Capital  string `json:"capital"`
	IsActive bool   `json:"is_active"`
}
