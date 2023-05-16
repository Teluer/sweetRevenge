package dto

import "time"

type FirstName struct {
	FirstName string `gorm:"primaryKey"`
	UsedLast  time.Time
	UsedTimes int
}
