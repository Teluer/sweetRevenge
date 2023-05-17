package dto

type LastName struct {
	LastName  string `gorm:"primaryKey"`
	UsedTimes int
}
