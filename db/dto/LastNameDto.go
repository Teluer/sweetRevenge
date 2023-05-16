package dto

import "time"

type LastName struct {
	LastName  string `gorm:"primaryKey"`
	UsedLast  time.Time
	UsedTimes int
}
