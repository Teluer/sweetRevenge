package dto

type ManualOrder struct {
	Name  string `gorm:"primaryKey;type:varchar(100)"`
	Phone string `gorm:"primaryKey;type:varchar(30)"`
}

func (ManualOrder) TableName() string {
	return "manual_order"
}
