package models

import "time"

// XAlarm 报警
type XAlarm struct {
	Id        int64     `gorm:"primary_key"`
	DeviceId  string    `json:"deviceId"`
	UUID      string    `json:"uuid"` // 报警ID
	StatusId  int64     `json:"statusId"`
	Status    uint8     `json:"status"`                          // 0-实时 1-补传
	Type      int       `json:"type"`                            // 报警类型
	StartTime time.Time `json:"startTime"`                       // 开始时间
	EndTime   time.Time `json:"endTime"`                         // 结束时间
	Data      string    `json:"data" gorm:"type:varchar(1024);"` // gps信息 json 字符串
}

// TableName 表名
func (s *XAlarm) TableName() string {
	return "xalarm"
}
