package impl

import (
	"encoding/json"
	"fmt"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/streadway/amqp"
)

// RabbitMQService implements the MessageBrokerService interface for RabbitMQ
type RabbitMQService struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	ExchangeName string
	ExchangeType string
}

// NewRabbitMQService creates a new RabbitMQ service
func NewRabbitMQService(conn *amqp.Connection, exchangeName string, exchangeType string) *RabbitMQService {
	channel, err := conn.Channel()
	exception.PanicLogging(err)

	// Declare the exchange
	err = channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	exception.PanicLogging(err)

	return &RabbitMQService{
		Connection:   conn,
		Channel:      channel,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
}

// PublishMessage publishes a message to a specified topic (routing key)
func (s *RabbitMQService) PublishMessage(routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error marshaling message: %s", err.Error()))
		return err
	}

	err = s.Channel.Publish(
		s.ExchangeName, // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		})

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error publishing message: %s", err.Error()))
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Published message to %s with routing key %s", s.ExchangeName, routingKey))
	return nil
}

// SubscribeToTopic subscribes to a topic (queue bound to exchange with routing key)
func (s *RabbitMQService) SubscribeToTopic(routingKey string, handler func([]byte) error) error {
	// Create a queue for this subscription
	queueName := fmt.Sprintf("%s.%s", s.ExchangeName, routingKey)
	queue, err := s.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error declaring queue: %s", err.Error()))
		return err
	}

	// Bind the queue to the exchange with the routing key
	err = s.Channel.QueueBind(
		queue.Name,     // queue name
		routingKey,     // routing key
		s.ExchangeName, // exchange
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error binding queue: %s", err.Error()))
		return err
	}

	// Create a consumer
	msgs, err := s.Channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error consuming queue: %s", err.Error()))
		return err
	}

	// Process messages in a goroutine
	go func() {
		for msg := range msgs {
			err := handler(msg.Body)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("Error handling message: %s", err.Error()))
				// Negative acknowledgement, message will be requeued
				msg.Nack(false, true)
			} else {
				// Acknowledge the message
				msg.Ack(false)
			}
		}
	}()

	logger.Logger.Info(fmt.Sprintf("Subscribed to %s with routing key %s", s.ExchangeName, routingKey))
	return nil
}

// Close closes the connection to RabbitMQ
func (s *RabbitMQService) Close() error {
	if err := s.Channel.Close(); err != nil {
		return err
	}
	return s.Connection.Close()
}

// PublishWithHeaders publishes a message with additional headers
func (s *RabbitMQService) PublishWithHeaders(routingKey string, message interface{}, headers amqp.Table) error {
	body, err := json.Marshal(message)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error marshaling message: %s", err.Error()))
		return err
	}

	err = s.Channel.Publish(
		s.ExchangeName, // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent, // Make message persistent
		})

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error publishing message with headers: %s", err.Error()))
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Published message with headers to %s with routing key %s", s.ExchangeName, routingKey))
	return nil
}

// CreateQueue creates a new queue
func (s *RabbitMQService) CreateQueue(queueName string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	queue, err := s.Channel.QueueDeclare(
		queueName,  // name
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		args,       // arguments
	)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error creating queue %s: %s", queueName, err.Error()))
		return queue, err
	}

	logger.Logger.Info(fmt.Sprintf("Created queue: %s", queueName))
	return queue, nil
}

// BindQueueToExchange binds a queue to an exchange with a routing key
func (s *RabbitMQService) BindQueueToExchange(queueName, routingKey, exchangeName string) error {
	err := s.Channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error binding queue %s to exchange %s: %s",
			queueName, exchangeName, err.Error()))
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Bound queue %s to exchange %s with routing key %s",
		queueName, exchangeName, routingKey))
	return nil
}
