package models

type Quote struct {
	BaseModel
	Quote  string `gorm:"column:quote"`
	Author string `gorm:"column:author"`
}
