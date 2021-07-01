package models

// 上下线信息

// XLink 定义
type XLink struct {
	Id            int64  `gorm:"primary_key"`
	Guid          string `json:"guid" gorm:"primary_key"`
	DeviceNo      string `json:"deviceNo" gorm:"type:varchar(24);"`
	RemoteAddress string `json:"remoteAddress"`                   // 设备网络地址
	OnTime        string `json:"onTime" gorm:"type:varchar(20);"` // 由设备上报
	OffTime       string `json:"offTime" gorm:"type:varchar(20);"`
	NetType       int    `json:"netType"`  // 网络类型
	Type          int    `json:"type"`     // 工作类型
	UpTraffic     int64  `json:"upFlow"`   // 上行流量
	DownTraffic   int64  `json:"downFlow"` // 下行流量
	Version       string `json:"vesion" gorm:"type:varchar(32);"`
}

// TableName 表名
func (s *XLink) TableName() string {
	return "t_xlink"
}
