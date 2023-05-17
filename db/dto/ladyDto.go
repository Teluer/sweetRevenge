package dto

type Lady struct {
	Phone     string `gorm:"primaryKey"`
	UsedTimes int
}

func (Lady) TableName() string {
	return "ladies"
}
