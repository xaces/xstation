package model

type DeviceOpt struct {
	Id         uint64   `json:"deviceId" gorm:"primary_key"`
	DeviceNo   string   `json:"deviceNo" gorm:"type:varchar(24);"`
	DeviceName string   `json:"deviceName" gorm:"type:varchar(20);"`
	ChnNumber  int      `json:"chnNumber"`
	ChnNames   JStrings `json:"chnNames"`
	Icon       string   `json:"icon" gorm:"type:varchar(64);"`
	Remark     string   `json:"remark" gorm:"size:500;"`
}

type Device struct {
	DeviceOpt
	Type       string `json:"type" gorm:"type:varchar(20);"`
	Guid       string `json:"guid" gorm:"type:varchar(64);"`
	Version    string `json:"version" gorm:"type:varchar(20);"`
	Online     bool   `json:"online"`
	DeviceTime string `json:"deviceTime" gorm:"type:varchar(20);"`
	TimeModel
}

// TableName 表名
func (s *Device) TableName() string {
	return "t_device"
}
