package dto

type ManualOrder struct {
	Name  string `gorm:"primaryKey;type:varchar(100)" json:"name"`
	Phone string `gorm:"primaryKey;type:varchar(30)"  json:"phone"`
}

func (ManualOrder) TableName() string {
	return "manual_order"
}
