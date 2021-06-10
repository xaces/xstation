package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/wlgd/xproto"
)

type Model struct {
	Id        uint64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
}

// jtime format json time field by myself
type jtime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t jtime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t jtime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *jtime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = jtime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

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
