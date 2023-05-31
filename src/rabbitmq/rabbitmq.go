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

type ManualOrder struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var rabbit struct {
	rabbitPool *pool.Pool
	queue      string
}

func InitializeRabbitMq(cfg config.RabbitConfig) {
	log.Info("Initializing rabbitmq connection")
	//wait for rabbitmq to initialize
	TestConnection(cfg.Host)

	rabbit.rabbitPool = openConnectionPool(cfg)
	rabbit.queue = cfg.QueueName

	// Declare queue
	ch := GetChannel()
	defer ch.Close()

	_, err := ch.QueueDeclare(
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

		log.Warn("Failed to connect to RabbitMQ: %v", err)
		time.Sleep(retryTimeout)
	}
	log.Error("Failed to connect to connect to RabbitMQ after %d attempts!", maxRetries)
	return false
}

func GetChannel() *amqp.Channel {
	conn, err := (*rabbit.rabbitPool).Get()
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

func ConsumeManualOrder() *ManualOrder {
	defer util.RecoverAndLog("RabbitMq")

	ch := GetChannel()
	defer ch.Close()
	messages, err := ch.Consume(
		rabbit.queue, // queue
		"orders",     // consumer
		true,         // auto-acknowledge messages
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Errorf("Failed to register a consumer: %v", err)
		panic(err)
	}

	message := <-messages
	log.Infof("Received message: %s", string(message.Body))

	var manualOrder ManualOrder
	err = json.Unmarshal(message.Body, &manualOrder)
	if err != nil {
		log.Errorf("Failed to unmarshal message: %v", err)
		// Handle unmarshal error
		return nil
	}
	return &manualOrder
}

func Publish(order *ManualOrder) error {
	// RabbitMQ connection URL
	ch := GetChannel()
	defer ch.Close()

	// Publish a message to the queue
	message, err := json.Marshal(*order)
	if err != nil {
		log.Error("Failed to marshal manual order!")
		return err
	}

	err = ch.Publish(
		"",           // exchange
		rabbit.queue, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		log.Error("Failed to publish a message: %v", err)
		return err
	} else {
		log.Infof("Message sent: %s", message)
	}
	return nil
}

func openConnectionPool(cfg config.RabbitConfig) *pool.Pool {
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
	return &connPool
}
