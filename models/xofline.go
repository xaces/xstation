package models

import "time"

// 上下线信息

// OFLine 定义
type XOFLine struct {
	Id            int64     `gorm:"primary_key"`
	Guid          string    `json:"guid" gorm:"primary_key"`
	DeviceId      string    `json:"deviceId"`
	RemoteAddress string    `json:"remoteAddress"` // 设备网络地址
	OnTime        time.Time `json:"onTime"`        // 由设备上报
	OffTime       time.Time `json:"offTime"`
	AccessType    int       `json:"access"`   // 网络类型
	Type          int       `json:"type"`     // 工作类型
	UpFlow        int64     `json:"upFlow"`   // 上行流量
	DownFlow      int64     `json:"downFlow"` // 下行流量
	Version       string    `json:"vesion"`
}

// TableName 表名
func (s *XOFLine) TableName() string {
	return "xofline"
}
