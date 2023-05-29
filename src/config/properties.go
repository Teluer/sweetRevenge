package config

import (
	"time"
)

type Config struct {
	TimeZone string `properties:"timezone"`

	OrdersRoutineCfg OrdersRoutineConfig `properties:"ordersroutine"`
	Rabbit           RabbitConfig        `properties:"rabbit"`

	SocksProxyAddress string `properties:"socks.proxy"`
	FirstNamesUrl     string `properties:"url.firstnames"`
	LastNamesUrl      string `properties:"url.lastnames"`
	DatabaseDsn       string `properties:"db.dsn"`

	LadiesCfg LadiesConfig `properties:"ladies"`
}

type OrdersRoutineConfig struct {
	SendOrdersMaxInterval time.Duration `properties:"send.interval.max"`
	StartDelay            bool          `properties:"start.delay"`
	DayStart              time.Duration `properties:"day.start"`
	DayEnd                time.Duration `properties:"day.end"`
	SendOrdersEnabled     bool          `properties:"orders.enabled""`
	OrdersCfg             OrdersConfig  `properties:"orders"`
}

type OrdersConfig struct {
	//loading string array as string to avoid space trimming by the properties lib
	PhonePrefixes    string   `properties:"phone.prefixes"`
	TargetBaselink   string   `properties:"target.base"`
	TargetOrderLink  string   `properties:"target.order"`
	TargetCategories []string `properties:"target.categories"`
}

type RabbitConfig struct {
	Host      string `properties:"host"`
	QueueName string `properties:"queue"`
}

type LadiesConfig struct {
	UpdateLadiesInterval time.Duration `properties:"update.interval"`
	LadiesBaseUrl        string        `properties:"base"`
	LadiesUrls           []string      `properties:"categories"`
}
