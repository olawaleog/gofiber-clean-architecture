package model

// QueuedEmailMessage represents an email message ready to be queued in RabbitMQ
type QueuedEmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	From    string `json:"from,omitempty"`
}
