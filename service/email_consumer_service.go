package service

// EmailConsumerService handles processing of queued email messages
type EmailConsumerService interface {
	// StartConsumer starts consuming email messages from the queue
	StartConsumer() error

	// ProcessEmail processes a single email message
	ProcessEmail(emailData []byte) error
}
