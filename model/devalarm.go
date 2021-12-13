package model

type DevAlarmOpt struct {
	DTU      string `json:"dtu"`
	DeviceNo string `json:"deviceNo" gorm:"type:varchar(24);"`
	Type     int    `json:"type"`                            // 类型
	Data     string `json:"data" gorm:"type:varchar(1024);"` // gps信息 json 字符串
}

// DevAlarm 报警
type DevAlarm struct {
	Guid string `json:"guid" gorm:"primary_key"`
	DevAlarmOpt
	Status    int        `json:"status"`                             // 0-开始 1--结束 2--报警中
	Flag      uint8      `json:"flag"`                               // 0-实时 1-补传
	StartTime string     `json:"startTime" gorm:"type:varchar(20);"` // 开始时间
	EndTime   string     `json:"endTime" gorm:"type:varchar(20);"`   // 结束时间
	StatusId  uint64     `json:"statusId"`
	DevStatus JDevStatus `json:"devStatus"`
}

// TableName 表名
func (s *DevAlarm) TableName() string {
	return "t_devalarm"
}

type dataType int

const (
	AlarmLinkUnknow       = dataType(0x00)
	AlarmLinkDev          = dataType(0x01)
	AlarmLinkFtpFile      = dataType(0x02)
	AlarmLinkDownloadFile = dataType(0x03)
)

// 报警关联信息
type DevAlarmLink struct {
	Id       uint64   `json:"id" gorm:"primary_key"`
	Guid     string   `json:"guid"`     // guid和DevAlarm.guid用来关联
	LinkType dataType `json:"linkType"` // 数据类型
	DevAlarmOpt
	DevStatus JDevStatus `json:"devStatus"`
}

// TableName 表名
func (s *DevAlarmLink) TableName() string {
	return "t_devalarmlink"
}
