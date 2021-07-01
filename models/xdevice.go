package models

type XDeviceOpt struct {
	DeviceNo   string `json:"deviceNo" gorm:"type:varchar(24);"`
	DeviceName string `json:"deviceName" gorm:"type:varchar(20);"`
	Icons      string `json:"icon" gorm:"type:varchar(64);"`
}

type XDevice struct {
	Model
	Guid    string `json:"guid" gorm:"type:varchar(64);"`
	Version string `json:"version" gorm:"type:varchar(20);"`
	Type    string `json:"type" gorm:"type:varchar(20);"`
	XDeviceOpt
}

// TableName 表名
func (s *XDevice) TableName() string {
	return "t_xdevice"
}
