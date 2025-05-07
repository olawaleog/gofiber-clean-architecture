package model

import "time"

type TransactionStatusModel struct {
	Data struct {
		Amount        int `json:"amount"`
		Authorization struct {
			AccountName       interface{} `json:"account_name"`
			AuthorizationCode string      `json:"authorization_code"`
			Bank              string      `json:"bank"`
			Bin               string      `json:"bin"`
			Brand             string      `json:"brand"`
			CardType          string      `json:"card_type"`
			Channel           string      `json:"channel"`
			CountryCode       string      `json:"country_code"`
			ExpMonth          string      `json:"exp_month"`
			ExpYear           string      `json:"exp_year"`
			Last4             string      `json:"last4"`
			MobileMoneyNumber string      `json:"mobile_money_number"`
			Reusable          bool        `json:"reusable"`
			Signature         interface{} `json:"signature"`
		} `json:"authorization"`
		Channel    string      `json:"channel"`
		Connect    interface{} `json:"connect"`
		CreatedAt  time.Time   `json:"createdAt"`
		CreatedAt1 time.Time   `json:"created_at"`
		Currency   string      `json:"currency"`
		Customer   struct {
			CustomerCode             string      `json:"customer_code"`
			Email                    string      `json:"email"`
			FirstName                interface{} `json:"first_name"`
			Id                       int         `json:"id"`
			InternationalFormatPhone interface{} `json:"international_format_phone"`
			LastName                 interface{} `json:"last_name"`
			Metadata                 interface{} `json:"metadata"`
			Phone                    interface{} `json:"phone"`
			RiskAction               string      `json:"risk_action"`
		} `json:"customer"`
		Domain          string      `json:"domain"`
		Fees            int         `json:"fees"`
		FeesBreakdown   interface{} `json:"fees_breakdown"`
		FeesSplit       interface{} `json:"fees_split"`
		GatewayResponse string      `json:"gateway_response"`
		Id              int64       `json:"id"`
		IpAddress       string      `json:"ip_address"`
		Log             struct {
			Attempts int `json:"attempts"`
			Errors   int `json:"errors"`
			History  []struct {
				Message string `json:"message"`
				Time    int    `json:"time"`
				Type    string `json:"type"`
			} `json:"history"`
			Input     []interface{} `json:"input"`
			Mobile    bool          `json:"mobile"`
			StartTime int           `json:"start_time"`
			Success   bool          `json:"success"`
			TimeSpent int           `json:"time_spent"`
		} `json:"log"`
		Message    interface{} `json:"message"`
		Metadata   int         `json:"metadata"`
		OrderId    interface{} `json:"order_id"`
		PaidAt     time.Time   `json:"paidAt"`
		PaidAt1    time.Time   `json:"paid_at"`
		Plan       interface{} `json:"plan"`
		PlanObject struct {
		} `json:"plan_object"`
		PosTransactionData interface{} `json:"pos_transaction_data"`
		ReceiptNumber      string      `json:"receipt_number"`
		Reference          string      `json:"reference"`
		RequestedAmount    int         `json:"requested_amount"`
		Source             interface{} `json:"source"`
		Split              struct {
		} `json:"split"`
		Status     string `json:"status"`
		Subaccount struct {
		} `json:"subaccount"`
		TransactionDate time.Time `json:"transaction_date"`
		TransactionId   uint      `json:"transaction_id"`
	} `json:"data"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
