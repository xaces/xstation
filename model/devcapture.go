package model

type DevCapture struct {
	Id        uint    `json:"id" gorm:"primary_key"`
	DeviceNo  string  `json:"deviceNo"`
	DTU       string  `json:"dtu" gorm:"type:datetime;default:null"`
	Channel   int     `json:"channel"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Speed     float32 `json:"speed"`
	Name      string  `json:"name"`
	Data      string  `json:"data"`
}

// TableName 表名
func (DevCapture) TableName() string {
	return "t_devcapture"
}
