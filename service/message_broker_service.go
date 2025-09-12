package service

// MessageBrokerService interface for messaging services
type MessageBrokerService interface {
	// PublishMessage publishes a message to a specified topic
	PublishMessage(topic string, message interface{}) error

	// SubscribeToTopic subscribes to a topic and processes messages with the handler
	SubscribeToTopic(topic string, handler func([]byte) error) error

	// Close closes the connection to the message broker
	Close() error
}
