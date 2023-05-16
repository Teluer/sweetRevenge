package dto

import "time"

type Lady struct {
	Phone     string `gorm:"primaryKey"`
	UsedLast  time.Time
	UsedTimes int
}

func (Lady) TableName() string {
	return "ladies"
}
