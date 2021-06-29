package models

import (
	"database/sql/driver"

	jsoniter "github.com/json-iterator/go"
	"github.com/wlgd/xproto"
)

// JGps
type JGps xproto.Gps

// Value insert
func (j JGps) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JGps) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

// JFloats float array
type JFloats []float32

// Value insert
func (j JFloats) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JFloats) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

// JMileage mileage
type JMileage xproto.Mileage

// Value insert
func (j JMileage) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JMileage) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

type JOil []xproto.Oil

// Value insert
func (j JOil) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JOil) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

type JModule xproto.Module

// Value insert
func (j JModule) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JModule) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

type JGsensor xproto.Gsensor

// Value insert
func (j JGsensor) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JGsensor) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

type JMobile xproto.Mobile

// Value insert
func (j JMobile) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JMobile) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

type JDisks []xproto.Disk

// Value insert
func (j JDisks) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JDisks) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

type JPeople xproto.People

// Value insert
func (j JPeople) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JPeople) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

type JObds []xproto.Obd

// Value insert
func (j JObds) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JObds) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), t)
}

// XStatus 状态数据
// gps {"longtitude:,latitude:,..."}
type XStatus struct {
	Id        int64    `json:"id" gorm:"primary_key"`
	DeviceId  uint64   `json:"deviceId"`
	DeviceNo  string   `json:"deviceNo" gorm:"type:varchar(24);"` // 时间
	Status    uint8    `json:"status"`                            // 0-实时 1-补传
	Acc       uint8    `json:"acc"`                               // acc
	DTU       string   `json:"dtu" gorm:"type:varchar(20);"`      // 时间
	Gps       JGps     `json:"gps" gorm:"type:varchar(128);"`     // gps信息 json 字符串
	Obds      JObds    `json:"obds"`                              // obd json 字符串
	Tempers   JFloats  `json:"tempers" gorm:"type:varchar(32);"`  // 温度 json 字符串
	Humiditys JFloats  `json:"humidity" gorm:"type:varchar(32);"` // 湿度 json 字符串
	Mileage   JMileage `json:"mileage" gorm:"type:varchar(64);"`  // 里程 json 字符串
	Oils      JOil     `json:"oils" gorm:"type:varchar(128);"`    // 油耗 json 字符串
	Module    JModule  `json:"module" gorm:"type:varchar(64);"`   // 模块状态 json 字符串
	Gsensor   JGsensor `json:"gsensor" gorm:"type:varchar(64);"`  // GSensor json 字符串
	Mobile    JMobile  `json:"mobile" gorm:"type:varchar(64);"`   // 移动网络 json 字符串
	Disks     JDisks   `json:"disks" gorm:"type:varchar(128);"`   // 磁盘 json 字符串
	People    JPeople  `json:"people" gorm:"type:varchar(32);"`   // 人数统计 json 字符串
	TableIdx  int      `json:"-" gorm:"-"`
}

// TableName 表名
func (s *XStatus) TableName() string {
	return "t_xstatus0"
}

type XStatus1 XStatus

// TableName 表名
func (s *XStatus1) TableName() string {
	return "t_xstatus1"
}

type XStatus2 XStatus

// TableName 表名
func (s *XStatus2) TableName() string {
	return "t_xstatus2"
}

type XStatus3 XStatus

// TableName 表名
func (s *XStatus3) TableName() string {
	return "t_xstatus3"
}

type XStatus4 XStatus

// TableName 表名
func (s *XStatus4) TableName() string {
	return "t_xstatus4"
}

var (
	KXStatusTabNumber int = 2
)
