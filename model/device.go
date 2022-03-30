package model

type DeviceBasic struct {
	ChlCount int    `json:"chlCount"`
	ChlNames string `json:"chlNames" gorm:"comment:通道别名,隔开;"`
	IoCount  int    `json:"ioCount"`
	IoNames  string `json:"ioNames" gorm:"comment:io别名,隔开;"`
}

type DeviceOpt struct {
	Id         uint64 `json:"deviceId" gorm:"primary_key"`
	DeviceNo   string `json:"deviceNo"`
	DeviceName string `json:"deviceName"`
	Icon       string `json:"icon"`
	DeviceBasic
	OrganizeId   uint64 `json:"organizeId"` // 分组Id
	OrganizeGuid string `json:"organizeGuid"`
	Details      string `json:"details"`
}

type Device struct {
	DeviceOpt
	Type           string `json:"type" gorm:"type:varchar(20);"`
	Guid           string `json:"guid" gorm:"type:varchar(64);"`
	Version        string `json:"version" gorm:"type:varchar(20);"`
	Online         bool   `json:"online"`
	LastOnlineTime string `json:"lastOnlineTime"`
	EffectiveTime  string `json:"effectiveTime"`
	ModelTime
}

// TableName 表名
func (s *Device) TableName() string {
	return "t_device"
}
