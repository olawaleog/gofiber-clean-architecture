package model

// QueuedSMSMessage represents an SMS message queued for sending through RabbitMQ
type QueuedSMSMessage struct {
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
	Message     string `json:"message"`
	UserID      uint   `json:"user_id,omitempty"` // Optional user ID to associate the message with
}
