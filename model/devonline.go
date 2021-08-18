package model

// 上线信息

// OnLine 定义
type DevOnline struct {
	Id            int64  `json:"id" gorm:"primary_key"`
	Guid          string `json:"guid" gorm:"primary_key"`
	DeviceNo      string `json:"deviceNo" gorm:"type:varchar(24);"`
	RemoteAddress string `json:"remoteAddress"`                   // 设备网络地址
	OnTime        string `json:"onTime" gorm:"type:varchar(20);"` // 由设备上报
	OffTime       string `json:"offTime" gorm:"type:varchar(20);"`
	NetType       int    `json:"netType"`     // 网络类型
	Type          int    `json:"type"`        // 工作类型
	UpTraffic     int64  `json:"upTraffic"`   // 上行流量
	DownTraffic   int64  `json:"downTraffic"` // 下行流量
}

// TableName 表名
func (s *DevOnline) TableName() string {
	return "t_devonline"
}