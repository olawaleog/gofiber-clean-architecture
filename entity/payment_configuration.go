package entity

import (
	"gorm.io/gorm"
)

type PaymentConfiguration struct {
	gorm.Model
	CountryCode string `gorm:"column:country_code;type:varchar(10);unique"`
	SecretKey   string `gorm:"column:secret_key;type:text"`
	PublicKey   string `gorm:"column:public_key;type:text"`
}

func (PaymentConfiguration) TableName() string {
	return "tb_payment_configurations"
}
