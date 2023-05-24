package config

import (
	"time"
)

type Config struct {
	OrdersRoutineCfg OrdersRoutineConfig
	Rabbit           RabbitConfig `properties:"Rabbit"`

	SocksProxyAddress string `properties:"SocksProxy"`
	FirstNamesUrl     string `properties:"FirstNamesUrl"`
	LastNamesUrl      string `properties:"LastNamesUrl"`

	LadiesCfg LadiesConfig
}

type OrdersRoutineConfig struct {
	SendOrdersMaxInterval time.Duration `properties:"SendOrdersMaxInterval"`
	DayStart              time.Duration `properties:"DayStart"`
	DayEnd                time.Duration `properties:"DayEnd"`
	TimeZone              string        `properties:"TimeZone"`
	SendOrdersEnabled     bool          `properties:"orders.enabled""`
	OrdersCfg             OrdersConfig
}

type OrdersConfig struct {
	PhonePrefixes    []string `properties:"PhonePrefixes"`
	TargetBaselink   string   `properties:"TargetBaselink"`
	TargetOrderLink  string   `properties:"TargetOrderLink"`
	TargetCategories []string `properties:"TargetCategories"`
}

type RabbitConfig struct {
	Host      string
	QueueName string
}

type LadiesConfig struct {
	UpdateLadiesInterval time.Duration `properties:"UpdateLadiesInterval"`
	LadiesBaseUrl        string        `properties:"LadiesBaseUrl"`
	LadiesUrls           []string      `properties:"LadiesUrls"`
}
