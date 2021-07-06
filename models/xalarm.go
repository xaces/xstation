package models

// XAlarm 报警
type XAlarm struct {
	Id       int64  `gorm:"primary_key"`
	Guid     string `json:"guid" gorm:"type:varchar(64);"`
	DeviceNo string `json:"deviceNo" gorm:"type:varchar(24);"`
	UUID      string `json:"uuid" gorm:"type:varchar(32);"` // 报警ID
	StatusId  int64  `json:"statusId"`
	Status    uint8  `json:"status"`                             // 0-实时 1-补传
	Type      int    `json:"type"`                               // 报警类型
	StartTime string `json:"startTime" gorm:"type:varchar(20);"` // 开始时间
	EndTime   string `json:"endTime" gorm:"type:varchar(20);"`   // 结束时间
	Data      string `json:"data" gorm:"type:varchar(1024);"`    // gps信息 json 字符串
}

// TableName 表名
func (s *XAlarm) TableName() string {
	return "t_xalarm"
}
