package model

// DevAlarm 报警
type DevAlarm struct {
	Guid      string     `json:"guid" gorm:"primary_key"`
	DeviceNo  string     `json:"deviceNo" gorm:"type:varchar(24);"`
	StatusId  uint64     `json:"statusId"`
	Status    JDevStatus `json:"status"`
	Flag      uint8      `json:"flag"`                               // 0-实时 1-补传
	Type      int        `json:"type"`                               // 报警类型
	StartTime string     `json:"startTime" gorm:"type:varchar(20);"` // 开始时间
	EndTime   string     `json:"endTime" gorm:"type:varchar(20);"`   // 结束时间
	Data      string     `json:"data" gorm:"type:varchar(1024);"`    // gps信息 json 字符串
}

// TableName 表名
func (s *DevAlarm) TableName() string {
	return "t_devalarm"
}
