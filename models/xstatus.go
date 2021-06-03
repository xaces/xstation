package models

import "time"

// XStatus 状态数据
// gps {"longtitude:,latitude:,..."}
type XStatus struct {
	Id        int64     `gorm:"primary_key"`
	DeviceId  string    `json:"deviceId"`
	Status    uint8     `json:"status"`   // 0-实时 1-补传
	DTU       time.Time `json:"dtu"`      // 时间
	Gps       string    `json:"gps"`      // gps信息 json 字符串
	Obds      string    `json:"obds"`     // obd json 字符串
	Tempers   string    `json:"tempers"`  // 温度 json 字符串
	Humiditys string    `json:"humidity"` // 湿度 json 字符串
	Mileage   string    `json:"mileage"`  // 里程 json 字符串
	Oils      string    `json:"oils"`     // 油耗 json 字符串
	Module    string    `json:"module"`   // 模块状态 json 字符串
	Gsensor   string    `json:"gsensor"`  // GSensor json 字符串
	Mobile    string    `json:"mobile"`   // 移动网络 json 字符串
	Disks     string    `json:"disks"`    // 磁盘 json 字符串
	People    string    `json:"people"`   // 人数统计 json 字符串
}

// TableName 表名
func (s *XStatus) TableName() string {
	return "xstatus"
}
