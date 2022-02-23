package model

type DeviceOpt struct {
	Id        uint64 `json:"id" gorm:"primary_key"`
	DeviceNo  string `json:"deviceNo"`
	Name      string `json:"name"`
	ChlsCount int    `json:"chlsCount"`
	ChlsName  string `json:"chlsName" gorm:"comment:通道别名已,隔开;"`
	Icon      string `json:"icon"`
	AutoFtp   bool   `json:"autoFtp"`
	Remark    string `json:"remark" gorm:"size:500;"`
}

type Device struct {
	DeviceOpt
	Type       string     `json:"type" gorm:"type:varchar(20);"`
	Guid       string     `json:"guid" gorm:"type:varchar(64);"`
	Version    string     `json:"version" gorm:"type:varchar(20);"`
	Online     bool       `json:"online"`
	LastStatus JDevStatus `json:"lastStatus" gorm:"comment:离线时状态;"`
	ModelTime
}

// TableName 表名
func (s *Device) TableName() string {
	return "t_device"
}
