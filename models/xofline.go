package models

// 上下线信息

// OFLine 定义
type XOFLine struct {
	Id            int64  `gorm:"primary_key"`
	Guid          string `json:"guid" gorm:"primary_key"`
	DeviceId      string `json:"deviceId" gorm:"type:varchar(12);"`
	RemoteAddress string `json:"remoteAddress"`                   // 设备网络地址
	OnTime        string `json:"onTime" gorm:"type:varchar(20);"` // 由设备上报
	OffTime       string `json:"offTime" gorm:"type:varchar(20);"`
	AccessType    int    `json:"access"`   // 网络类型
	Type          int    `json:"type"`     // 工作类型
	UpFlow        int64  `json:"upFlow"`   // 上行流量
	DownFlow      int64  `json:"downFlow"` // 下行流量
	Version       string `json:"vesion" gorm:"type:varchar(32);"`
}

// TableName 表名
func (s *XOFLine) TableName() string {
	return "t_xofline"
}
