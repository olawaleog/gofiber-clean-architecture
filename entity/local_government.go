package entity

import "gorm.io/gorm"

type LocalGovernmentArea struct {
	gorm.Model
	Name    string `gorm:"column:name;type:varchar(100)"`
	Capital string `gorm:"column:capital;type:varchar(100)"`
}

func (LocalGovernmentArea) TableName() string {
	return "tb_local_government_areas"
}
