package dto

type FirstName struct {
	FirstName string `gorm:"primaryKey"`
	UsedTimes int
}
