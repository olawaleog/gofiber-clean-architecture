package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

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
	mu           sync.Mutex
	closeChan    chan *amqp.Error
	connClose    chan *amqp.Error
}

// NewRabbitMQService creates a new RabbitMQ service
func NewRabbitMQService(conn *amqp.Connection, exchangeName string, exchangeType string) *RabbitMQService {
	if conn == nil {
		exception.PanicLogging(errors.New("nil rabbitmq connection"))
	}

	s := &RabbitMQService{
		Connection:   conn,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
		closeChan:    make(chan *amqp.Error),
		connClose:    make(chan *amqp.Error),
	}

	if err := s.ensureChannelAndExchange(); err != nil {
		exception.PanicLogging(err)
	}

	// monitor connection close
	go func() {
		notify := s.Connection.NotifyClose(make(chan *amqp.Error))
		for err := range notify {
			if err != nil {
				logger.Logger.Warn(fmt.Sprintf("rabbitmq connection closed: %s", err.Error()))
			}
			// signal and keep loop ended
			close(s.connClose)
			return
		}
	}()

	return s
}

func (s *RabbitMQService) ensureChannelAndExchange() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Connection == nil || s.Connection.IsClosed() {
		return errors.New("rabbitmq connection is closed")
	}

	// If channel is already present and not nil, return
	if s.Channel != nil {
		return nil
	}

	ch, err := s.Connection.Channel()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error creating channel: %s", err.Error()))
		return err
	}

	// Declare the exchange (idempotent)
	err = ch.ExchangeDeclare(
		s.ExchangeName, // name
		s.ExchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		_ = ch.Close()
		logger.Logger.Error(fmt.Sprintf("Error declaring exchange: %s", err.Error()))
		return err
	}

	s.Channel = ch

	// watch channel close and clear s.Channel when closed
	notify := s.Channel.NotifyClose(make(chan *amqp.Error))
	go func() {
		if e, ok := <-notify; ok && e != nil {
			logger.Logger.Warn(fmt.Sprintf("rabbitmq channel closed: %s", e.Error()))
		} else {
			logger.Logger.Warn("rabbitmq channel closed")
		}
		s.mu.Lock()
		// mark channel nil so next publish re-creates it
		s.Channel = nil
		s.mu.Unlock()
	}()

	logger.Logger.Info(fmt.Sprintf("Declared exchange and created channel: %s", s.ExchangeName))
	return nil
}

// PublishMessage publishes a message to a specified topic (routing key)
func (s *RabbitMQService) PublishMessage(routingKey string, message interface{}) error {
	if err := s.ensureChannelAndExchange(); err != nil {
		logger.Logger.Error(fmt.Sprintf("Cannot publish, ensureChannel error: %s", err.Error()))
		return err
	}

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
		// if channel/connection closed, clear channel so next attempt recreates it
		s.mu.Lock()
		_ = s.Channel // attempt to leave it for Close handling, but mark nil
		s.Channel = nil
		s.mu.Unlock()
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Published message to %s with routing key %s", s.ExchangeName, routingKey))
	return nil
}

// SubscribeToTopic subscribes to a topic (queue bound to exchange with routing key)
func (s *RabbitMQService) SubscribeToTopic(routingKey string, handler func([]byte) error) error {
	if err := s.ensureChannelAndExchange(); err != nil {
		logger.Logger.Error(fmt.Sprintf("Cannot subscribe, ensureChannel error: %s", err.Error()))
		return err
	}

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
		// When msgs channel closes, mark channel nil so it can be recreated
		s.mu.Lock()
		s.Channel = nil
		s.mu.Unlock()
	}()

	logger.Logger.Info(fmt.Sprintf("Subscribed to %s with routing key %s", s.ExchangeName, routingKey))
	return nil
}

// Close closes the connection to RabbitMQ
func (s *RabbitMQService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	if s.Channel != nil {
		if e := s.Channel.Close(); e != nil {
			err = e
		}
		s.Channel = nil
	}
	if s.Connection != nil && !s.Connection.IsClosed() {
		if e := s.Connection.Close(); e != nil && err == nil {
			err = e
		}
	}
	return err
}

// PublishWithHeaders publishes a message with additional headers
func (s *RabbitMQService) PublishWithHeaders(routingKey string, message interface{}, headers amqp.Table) error {
	if err := s.ensureChannelAndExchange(); err != nil {
		logger.Logger.Error(fmt.Sprintf("Cannot publish with headers, ensureChannel error: %s", err.Error()))
		return err
	}

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
		s.mu.Lock()
		s.Channel = nil
		s.mu.Unlock()
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Published message with headers to %s with routing key %s", s.ExchangeName, routingKey))
	return nil
}

// CreateQueue creates a new queue
func (s *RabbitMQService) CreateQueue(queueName string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	if err := s.ensureChannelAndExchange(); err != nil {
		return amqp.Queue{}, err
	}

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
	if err := s.ensureChannelAndExchange(); err != nil {
		return err
	}

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
