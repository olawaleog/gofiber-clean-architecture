package entity

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	UserId float64 `gorm:"column:user_id;type:int"`
	User   User    `gorm:"foreignkey:UserId"`

	Amount float64 `gorm:"column:amount;type:numeric(10,2)"`

	WaterCost   float64 `gorm:"column:refinery_amount;type:numeric(10,2)"`
	DeliveryFee float64 `gorm:"column:delivery_amount;type:numeric(10,2)"`
	TotalAmount float64 `gorm:"column:total_amount;type:numeric(10,2)"`

	PhoneNumber string    `gorm:"column:phone_number;type:varchar(15)"`
	Email       string    `gorm:"column:email;type:varchar(255)"`
	Provider    string    `gorm:"column:provider;type:varchar(50)"`
	Status      string    `gorm:"column:status;type:varchar(20)"`
	PaymentID   string    `gorm:"column:payment_id;type:varchar(255)"`
	PaymentType string    `gorm:"column:payment_type;type:varchar(20)"`
	CompletedAt time.Time `gorm:"column:completed_at;type:timestamp"`
	Currency    string    `gorm:"column:currency;type:varchar(50)"`
	RawResponse string    `gorm:"column:raw_response;type:text"`
	Reference   string    `gorm:"column:reference;type:varchar(255)"`
	Scheme      string    `gorm:"column:scheme;type:varchar(255)"`
	RawRequest  string    `gorm:"column:raw_request;type:text"`
	//Order       Order     `gorm:"foreignkey:TransactionId"`
}

func (Transaction) TableName() string {
	return "tb_transactions"
}
