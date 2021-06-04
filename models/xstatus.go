package models

// XStatus 状态数据
// gps {"longtitude:,latitude:,..."}
type XStatus struct {
	Id        int64  `gorm:"primary_key"`
	DeviceId  string `json:"deviceId" gorm:"type:varchar(12);"`
	Status    uint8  `json:"status"`                            // 0-实时 1-补传
	DTU       string `json:"dtu" gorm:"type:varchar(20);"`      // 时间
	Gps       string `json:"gps" gorm:"type:varchar(128);"`     // gps信息 json 字符串
	Obds      string `json:"obds"`                              // obd json 字符串
	Tempers   string `json:"tempers" gorm:"type:varchar(32);"`  // 温度 json 字符串
	Humiditys string `json:"humidity" gorm:"type:varchar(32);"` // 湿度 json 字符串
	Mileage   string `json:"mileage" gorm:"type:varchar(32);"`  // 里程 json 字符串
	Oils      string `json:"oils" gorm:"type:varchar(128);"`    // 油耗 json 字符串
	Module    string `json:"module" gorm:"type:varchar(64);"`   // 模块状态 json 字符串
	Gsensor   string `json:"gsensor" gorm:"type:varchar(64);"`  // GSensor json 字符串
	Mobile    string `json:"mobile" gorm:"type:varchar(64);"`   // 移动网络 json 字符串
	Disks     string `json:"disks" gorm:"type:varchar(128);"`   // 磁盘 json 字符串
	People    string `json:"people" gorm:"type:varchar(32);"`   // 人数统计 json 字符串
}

// TableName 表名
func (s *XStatus) TableName() string {
	return "t_xstatus"
}
