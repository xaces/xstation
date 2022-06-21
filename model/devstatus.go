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

type JParam1 struct {
	Obds      []xproto.Obd  `json:"obds"`
	Tempers   []float32     `json:"temps"`
	Humiditys []float32     `json:"humis"`
	Module    xproto.Module `json:"mod"`   // 模块状态 json 字符串
	Disks     []xproto.Disk `json:"disks"` // 磁盘 json 字符串
}

// Value insert
func (j JParam1) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (j *JParam1) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), j)
}

type JParam2 struct {
	Gsensor xproto.Gsensor `json:"gs"`   // GSensor json 字符串
	People  xproto.People  `json:"peop"` // 人数统计 json 字符串
	Vols    []float32      `json:"vols"`
}

// Value insert
func (j JParam2) Value() (driver.Value, error) {
	return jsoniter.Marshal(&j)
}

// Scan valueof
func (j *JParam2) Scan(v interface{}) error {
	return jsoniter.Unmarshal(v.([]byte), j)
}

// Status 状态数据
type DevStatus struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	DeviceID uint      `json:"deviceId" gorm:"index:idx_status"`
	Flag     uint8     `json:"flag"`                                                  // 0-实时 1-补传 2-报警开始 3-报警结束
	Acc      uint8     `json:"acc"`                                                   // acc
	DTU      string    `json:"dtu" gorm:"type:datetime;primary_key;index:idx_status"` // 时间
	Location JLocation `json:"loc" gorm:"type:varchar(128);"`                         // location json 字符串
	Mileage  JMileage  `json:"mile" gorm:"type:varchar(32)"`                          // 里程 json 字符串
	Oils     JOil      `json:"oils"  gorm:"type:varchar(32)"`                         // 油耗 json 字符串
	P1       JParam1   `json:"p1" gorm:"type:varchar(512)"`
	P2       JParam2   `json:"p2" gorm:"type:varchar(512)"`
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
	gStatusTabCount uint = 2
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
