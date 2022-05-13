package model

type DeviceOpt struct {
	ID            uint   `json:"deviceId" gorm:"primary_key"`
	No            string `json:"deviceNo"`
	Name          string `json:"deviceName"`
	Icon          string `json:"icon"`
	ChlCount      int    `json:"chlCount"`
	ChlNames      string `json:"chlNames" gorm:"comment:通道别名,隔开;"`
	IoCount       int    `json:"ioCount"`
	IoNames       string `json:"ioNames" gorm:"comment:io别名,隔开;"`
	EffectiveTime string `json:"effectiveTime"`
	Details       string `json:"details"`
}

type Device struct {
	DeviceOpt
	Type           string `json:"type" gorm:"type:varchar(20);"`
	Version        string `json:"version" gorm:"type:varchar(20);"`
	Online         bool   `json:"online"`
	LastOnlineTime string `json:"lastOnlineTime" gorm:"type:datetime;default:null"`
	CreatedAt      jtime  `json:"createdAt"`
	UpdatedAt      jtime  `json:"updatedAt"`
}

// TableName 表名
func (DeviceOpt) TableName() string {
	return "t_device"
}
