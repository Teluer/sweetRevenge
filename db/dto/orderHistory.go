package dto

import "time"

type OrderHistory struct {
	Name          string
	Phone         string
	ItemId        string
	OrderDateTime time.Time
}

func (OrderHistory) TableName() string {
	return "order_history"
}
