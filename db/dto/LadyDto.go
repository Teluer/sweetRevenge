package dto

import "time"

type Lady struct {
	phone     string `gorm:"primaryKey"`
	usedLast  time.Time
	usedTimes int
}

func (Lady) TableName() string {
	return "ladies"
}
