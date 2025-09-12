package configuration

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/streadway/amqp"
)

// NewRabbitMQ creates a new RabbitMQ connection
func NewRabbitMQ(config Config) *amqp.Connection {
	connectionString := config.Get("RABBITMQ_URL")
	if connectionString == "" {
		// Fallback to constructing URL from parts
		user := config.Get("RABBITMQ_USER")
		password := config.Get("RABBITMQ_PASSWORD")
		host := config.Get("RABBITMQ_HOST")
		port := config.Get("RABBITMQ_PORT")
		vhost := config.Get("RABBITMQ_VHOST")

		// Default vhost to "/"
		if vhost == "" {
			vhost = "/"
		}

		// Construct the connection URL
		connectionString = "amqp://" + user + ":" + password + "@" + host + ":" + port + "/" + vhost
	}

	conn, err := amqp.Dial(connectionString)
	exception.PanicLogging(err)

	return conn
}

// CreateChannel creates a new channel from a RabbitMQ connection
func CreateChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	exception.PanicLogging(err)

	return ch
}

// DeclareQueue declares a queue to use
func DeclareQueue(ch *amqp.Channel, queueName string) amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	exception.PanicLogging(err)

	return q
}

// DeclareExchange declares an exchange
func DeclareExchange(ch *amqp.Channel, exchangeName string, exchangeType string) {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	exception.PanicLogging(err)
}

// BindQueue binds a queue to an exchange
func BindQueue(ch *amqp.Channel, queueName string, routingKey string, exchangeName string) {
	err := ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	exception.PanicLogging(err)
}
