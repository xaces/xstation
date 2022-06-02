package model

import (
	"database/sql/driver"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/wlgd/xproto"
)

// JLocation
type JLocation xproto.Location

// Value insert
func (j JLocation) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (t *JLocation) Scan(v interface{}) error {
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

// Status 状态数据
type DevStatus struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	DeviceID uint   `json:"deviceId"`
	DeviceNo string `json:"deviceNo" gorm:"index:idx_status"`
	Flag     uint8  `json:"flag"`                                                  // 0-实时 1-补传 2-报警开始 3-报警结束
	Acc      uint8  `json:"acc"`                                                   // acc
	DTU      string `json:"dtu" gorm:"type:datetime;primary_key;index:idx_status"` // 时间

	Location  JLocation `json:"location" gorm:"type:varchar(128);"` // location json 字符串
	Obds      JObds     `json:"obds"`                               // obd json 字符串
	Tempers   JFloats   `json:"tempers"`                            // 温度 json 字符串
	Humiditys JFloats   `json:"humidity"`                           // 湿度 json 字符串
	Mileage   JMileage  `json:"mileage"`                            // 里程 json 字符串
	Oils      JOil      `json:"oils"`                               // 油耗 json 字符串
	Module    JModule   `json:"module"`                             // 模块状态 json 字符串
	Gsensor   JGsensor  `json:"gsensor"`                            // GSensor json 字符串
	Mobile    JMobile   `json:"mobile"`                             // 移动网络 json 字符串
	Disks     JDisks    `json:"disks"`                              // 磁盘 json 字符串
	People    JPeople   `json:"people"`                             // 人数统计 json 字符串
	Vols      JFloats   `json:"vols"`
}

// Value insert
func (j DevStatus) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (j *DevStatus) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), j)
}

// 分库
var (
	gStatusTabCount uint = 5
)

// TableName 表名
func (o DevStatus) TableName() string {
	return fmt.Sprintf("t_devstatus_%0*d", 1, o.DeviceID%gStatusTabCount)
}

func (DevStatus) TableNameOf(deviceID uint) string {
	o := DevStatus{DeviceID: deviceID}
	return o.TableName()
}

func (DevStatus) TableCount() uint {
	return gStatusTabCount
}
