package rabbitmq

import (
	"github.com/silenceper/pool"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sweetRevenge/src/config"
	"time"
)

func InitializeRabbitMq(cfg config.RabbitConfig) *pool.Pool {
	//wait for rabbitmq to initialize
	TestConnection(cfg.Host)

	poolConfig := &pool.Config{
		InitialCap: 1,
		MaxCap:     10,
		MaxIdle:    5,
		Factory: func() (interface{}, error) {
			return amqp.Dial(cfg.Host)
		},
		Close: func(v interface{}) error {
			conn := v.(amqp.Connection)
			return conn.Close()
		},
	}
	connPool, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ connection pool: %v", err)
	}

	// Declare queue
	ch := GetChannel(&connPool)
	defer ch.Close()

	_, err = ch.QueueDeclare(
		cfg.QueueName, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare manual orders queue: %v", err)
	}

	return &connPool
}

func TestConnection(url string) bool {
	const maxRetries = 5
	const retryTimeout = time.Second * 5

	for i := 1; i <= maxRetries; i++ {
		log.Info("Connecting to RabbitMQ (attempt %d/%d)...", i, maxRetries)
		conn, err := amqp.Dial(url)
		if err == nil {
			conn.Close()
			return true
		}

		log.Printf("Failed to connect to RabbitMQ: %v", err)
		time.Sleep(retryTimeout)
	}
	return false
}

func GetChannel(connPool *pool.Pool) *amqp.Channel {
	conn, err := (*connPool).Get()
	if err != nil {
		log.Panic("Failed to acquire connection from pool:", err)
		panic(err)
	}

	// Convert the connection to the appropriate type (amqp.Connection)
	rabbitMQConn := conn.(*amqp.Connection)

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Panic("Failed to open a channel:", err)
		panic(err)
	}

	return ch
}

//TODO: implement
func ConsumeManualOrder(connPool *pool.Pool, queue string) {
	ch := GetChannel(connPool)
	messages, err := ch.Consume(
		queue,    // queue
		"orders", // consumer
		true,     // auto-acknowledge messages
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	message := <-messages
	body := string(message.Body)
	log.Printf("Received message: %s", body)
}

//TODO: implement
func Publish(connPool *pool.Pool, queue string) {
	// RabbitMQ connection URL
	ch := GetChannel(connPool)
	defer ch.Close()

	// Publish a message to the queue
	message := "Hello, RabbitMQ!"
	err := ch.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Error("Failed to publish a message: %v", err)
	} else {
		log.Infof("Message sent: %s", message)
	}
}
