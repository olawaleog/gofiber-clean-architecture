package entity

import "gorm.io/gorm"

type PaymentMethod struct {
	gorm.Model
	Name        string  `gorm:"column:name;type:varchar(100)"`
	Type        string  `gorm:"column:type;type:varchar(100)"`
	Scheme      string  `gorm:"column:scheme;type:varchar(100)"`
	Description string  `gorm:"column:description;type:text"`
	RawData     string  `gorm:"column:raw_data;type:text"`
	UserID      float64 `gorm:"column:user_id;type:numeric"`
}

func (PaymentMethod) TableName() string {
	return "tb_payment_methods"
}
