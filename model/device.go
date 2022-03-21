package model

type DeviceOpt struct {
	Id           uint64 `json:"deviceId" gorm:"primary_key"`
	DeviceNo     string `json:"deviceNo"`
	DeviceName   string `json:"deviceName"`
	ChlsCount    int    `json:"chlsCount"`
	ChlsName     string `json:"chlsName" gorm:"comment:通道别名已,隔开;"`
	Icon         string `json:"icon"`
	AutoFtp      bool   `json:"autoFtp"`
	OrganizeGuid string `json:"organizeGuid"` // 组织GUid
	OrganizeId   uint64 `json:"organizeId"`   // 分组Id
	Details      string `json:"details"`
}

type Device struct {
	DeviceOpt
	Type       string      `json:"type" gorm:"type:varchar(20);"`
	Guid       string      `json:"guid" gorm:"type:varchar(64);"`
	Version    string      `json:"version" gorm:"type:varchar(20);"`
	Online     bool        `json:"online"`
	ModelTime
}

// TableName 表名
func (s *Device) TableName() string {
	return "t_device"
}
