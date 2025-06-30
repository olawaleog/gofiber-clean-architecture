package model

import (
	"github.com/google/uuid"
	"time"
)

//generate TransactionModel form entity.Transaction

type TransactionModel struct {
	ID              uint         `json:"id"`
	Email           string       `json:"email"`
	PhoneNumber     string       `json:"phone_number"`
	Amount          float64      `json:"amount"`
	Currency        string       `json:"currency"`
	UserId          uint         `json:"user_id"`
	PaymentID       string       `json:"payment_id"`
	Provider        string       `json:"provider"`
	PaymentType     string       `json:"payment_type"`
	Status          string       `json:"status"`
	Reference       string       `json:"reference"`
	RawRequest      string       `json:"raw_request"`
	RawResponse     string       `json:"raw_response"`
	WaterCost       float64      `json:"water_cost"`
	DeliveryFee     float64      `json:"delivery_fee"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       string       `json:"updated_at"`
	DeliveryAddress AddressModel `json:"delivery_address_detail"`
}
type TransactionCreateUpdateModel struct {
	Id                 string                               `json:"id"`
	TotalPrice         int64                                `json:"total_price"`
	TransactionDetails []TransactionDetailCreateUpdateModel `json:"transaction_details"`
}

type TransactionDetailModel struct {
	Id            string `json:"id"`
	SubTotalPrice int64  `json:"sub_total_price" validate:"required"`
	Price         int64  `json:"price" validate:"required"`
	Quantity      int32  `json:"quantity" validate:"required"`
	Product       ProductModel
}

type TransactionDetailCreateUpdateModel struct {
	Id            string    `json:"id"`
	SubTotalPrice int64     `json:"sub_total_price" validate:"required"`
	Price         int64     `json:"price" validate:"required"`
	Quantity      int32     `json:"quantity" validate:"required"`
	ProductId     uuid.UUID `json:"product_id" validate:"required"`
	Product       ProductModel
}
