package rabbitmq

import (
	"encoding/json"
	"github.com/silenceper/pool"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sweetRevenge/src/config"
	"sweetRevenge/src/util"
	"time"
)

var rabbitPool *pool.Pool

func InitializeRabbitMq(cfg config.RabbitConfig) {
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
		log.Error("Failed to create RabbitMQ connection pool: %v", err)
	}
	rabbitPool = &connPool

	// Declare queue
	ch := GetChannel()
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
		log.Error("Failed to declare manual orders queue: %v", err)
	}
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

		log.Errorf("Failed to connect to RabbitMQ: %v", err)
		time.Sleep(retryTimeout)
	}
	return false
}

func GetChannel() *amqp.Channel {
	conn, err := (*rabbitPool).Get()
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

func ConsumeManualOrder(queue string) *ManualOrder {
	defer util.RecoverAndLogError("RabbitMq")

	ch := GetChannel()
	defer ch.Close()
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
		log.Error("Failed to register a consumer: %v", err)
		panic(err)
	}

	message := <-messages
	log.Printf("Received message: %s", string(message.Body))

	var manualOrder ManualOrder
	err = json.Unmarshal(message.Body, &manualOrder)
	if err != nil {
		log.Error("Failed to unmarshal message: %v", err)
		// Handle unmarshal error
		return nil
	}
	return &manualOrder
}

// TODO: implement this when there is admin panel
//func Publish(queue string) {
//	// RabbitMQ connection URL
//	ch := GetChannel()
//	defer ch.Close()
//
//	// Publish a message to the queue
//	message := "Hello, RabbitMQ!"
//	err := ch.Publish(
//		"",    // exchange
//		queue, // routing key
//		false, // mandatory
//		false, // immediate
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        []byte(message),
//		})
//	if err != nil {
//		log.Error("Failed to publish a message: %v", err)
//	} else {
//		log.Infof("Message sent: %s", message)
//	}
//}
