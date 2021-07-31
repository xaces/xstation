package model

type DeviceOpt struct {
	DeviceNo   string `json:"deviceNo" gorm:"type:varchar(24);"`
	DeviceName string `json:"deviceName" gorm:"type:varchar(20);"`
	Icons      string `json:"icon" gorm:"type:varchar(64);"`
	Type       string `json:"type" gorm:"type:varchar(20);"`
}

type Device struct {
	Model
	Guid string `json:"guid" gorm:"type:varchar(64);"`
	DeviceOpt
	Version string `json:"version" gorm:"type:varchar(20);"`
}

// TableName 表名
func (s *Device) TableName() string {
	return "t_device"
}
