package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	Id        uint64         `json:"id" gorm:"primary_key"`
	CreatedAt jtime          `json:"createdAt"`
	UpdatedAt jtime          `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
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