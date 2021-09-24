package model

import (
	"database/sql/driver"
	"encoding/json"
)

type JDevStatus DevStatus

// Value insert
func (j JDevStatus) Value() (driver.Value, error) {
	return json.Marshal(&j)
}

// Scan valueof
func (t *JDevStatus) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), t)
}

type DeviceOpt struct {
	Id        uint64   `json:"id" gorm:"primary_key"`
	No        string   `json:"no" gorm:"type:varchar(24);"`
	Name      string   `json:"name" gorm:"type:varchar(20);"`
	ChlNumber int      `json:"chlNumber"`
	ChlNames  JStrings `json:"chlNames"`
	Icon      string   `json:"icon" gorm:"type:varchar(64);"`
	Remark    string   `json:"remark" gorm:"size:500;"`
}

type Device struct {
	DeviceOpt
	Type       string     `json:"type" gorm:"type:varchar(20);"`
	Guid       string     `json:"guid" gorm:"type:varchar(64);"`
	Version    string     `json:"version" gorm:"type:varchar(20);"`
	Online     bool       `json:"online"`
	LastStatus JDevStatus `json:"lastStatus" gorm:"comment:离线时状态;"`
	TimeModel
}

// TableName 表名
func (s *Device) TableName() string {
	return "t_device"
}
