package entity

import "gorm.io/gorm"

type Setting struct {
	gorm.Model
	Key   string `gorm:"column:key;type:varchar(50);unique_index"`
	Value string `gorm:"column:value;type:text"`
}

func (Setting) TableName() string {
	return "tb_settings"
}
