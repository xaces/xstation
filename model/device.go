package model

type DeviceBasic struct {
	ChlsCount int    `json:"chlsCount"`
	ChlsName  string `json:"chlsName" gorm:"comment:通道别名,隔开;"`
	IosCount  int    `json:"iosCount"`
	IosName   string `json:"iosName" gorm:"comment:io别名,隔开;"`
}

type DeviceOpt struct {
	Id         uint64 `json:"deviceId" gorm:"primary_key"`
	DeviceNo   string `json:"deviceNo"`
	DeviceName string `json:"deviceName"`
	Icon       string `json:"icon"`
	DeviceBasic
	OrganizeGuid string `json:"organizeGuid"` // 组织GUid
	OrganizeId   uint64 `json:"organizeId"`   // 分组Id
	Details      string `json:"details"`
}

type Device struct {
	DeviceOpt
	Type           string `json:"type" gorm:"type:varchar(20);"`
	Guid           string `json:"guid" gorm:"type:varchar(64);"`
	Version        string `json:"version" gorm:"type:varchar(20);"`
	Online         bool   `json:"online"`
	LastOnlineTime string `json:"lastOnlineTime"`
	ValidTime      jtime  `json:"validTime"`
	ModelTime
}

// TableName 表名
func (s *Device) TableName() string {
	return "t_device"
}
