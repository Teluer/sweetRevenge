package config

import (
	"time"
)

type Config struct {
	OrdersRoutineCfg OrdersRoutineConfig

	SocksProxyAddress string `properties:"SocksProxy"`

	FirstNamesUrl string `properties:"FirstNamesUrl"`
	LastNamesUrl  string `properties:"LastNamesUrl"`

	LadiesCfg LadiesConfig
}

type OrdersRoutineConfig struct {
	SendOrdersBaseInterval      time.Duration `properties:"SendOrdersBaseInterval"`
	SendOrdersIntervalVariation float32       `properties:"SendOrdersIntervalVariation"`
	DayStart                    time.Duration `properties:"DayStart"`
	DayEnd                      time.Duration `properties:"DayEnd"`
	OrdersCfg                   OrdersConfig
}

type OrdersConfig struct {
	PhonePrefixes    []string `properties:"PhonePrefixes"`
	TargetBaselink   string   `properties:"TargetBaselink"`
	TargetOrderLink  string   `properties:"TargetOrderLink"`
	TargetCategories []string `properties:"TargetCategories"`
}

type LadiesConfig struct {
	UpdateLadiesInterval time.Duration `properties:"UpdateLadiesInterval"`
	LadiesBaseUrl        string        `properties:"LadiesBaseUrl"`
	LadiesUrls           []string      `properties:"LadiesUrls"`
}
