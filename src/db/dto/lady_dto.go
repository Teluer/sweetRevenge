package dto

type Lady struct {
	Phone     string `gorm:"primaryKey;type:varchar(15)"`
	UsedTimes int
}

func (Lady) TableName() string {
	return "ladies"
}
