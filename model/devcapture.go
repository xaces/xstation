package model

type DevCapture struct {
	Id        uint64   `json:"id" gorm:"primary_key"`
	DeviceNo  string   `json:"deviceNo"`
	DTU       string   `json:"dtu"`
	Channel   int      `json:"channel"`
	Latitude  float32  `json:"latitude"`
	Longitude float32  `json:"longitude"`
	Speed     float32  `json:"speed"`
	Name      string   `json:"name"`
	Data      JReserve `json:"data"`
}

// TableName 表名
func (s *DevCapture) TableName() string {
	return "t_devcapture"
}
