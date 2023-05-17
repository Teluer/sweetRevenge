package dto

type ManualOrder struct {
	Name   string
	Phone  string
	Target string
}

func (ManualOrder) TableName() string {
	return "manual_order"
}
