package dto

import "time"

type Item struct {
	id        string `gorm:"primaryKey"`
	link      string
	category  string
	usedLast  time.Time
	usedTimes int
}

func (Item) TableName() string {
	return "goods"
}
