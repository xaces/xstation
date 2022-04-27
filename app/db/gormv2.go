package db

import (
	"fmt"
	"time"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

type Option struct {
	Name    string
	Address string
}

const (
	timeFormat = "2006-01-02 15:04:05"
	partFormat = "20060102 150405"
)

func Partition(table string, day int) {
	t := time.Now()
	where := t.Format(timeFormat)[:10]
	addp := "p" + t.Format(partFormat)[:8]
	delp := "p" + t.AddDate(0, 0, -1*day).Format(partFormat)[:8]
	fmt.Printf("ALTER %s PARTITION DROP %s ADD %s VALUES LESS THAN TO_DAYS %s", table, delp, addp, where)
	// orm.DB().Exec(fmt.Sprintf("ALTER TABLE %s DROP PARTITION %s", table, delp))
	// orm.DB().Exec(fmt.Sprintf("ALTER TABLE %s ADD PARTITION( PARTITION %s VALUES LESS THAN (TO_DAYS('%s')));", table, addp, where))
}

func Run(o *Option) error {
	db, err := orm.NewGormV2(o.Name, o.Address)
	if err != nil {
		return err
	}
	db.AutoMigrate(
		&model.DevOnline{},
		&model.DevAlarm{},
		&model.DevAlarmDetails{},
		&model.DevAlarmFile{},
		&model.DevCapture{},
		&model.DevStatus{},
		&model.DevStatus1{},
	)
	Partition("t_devalarm", 10)
	return nil
}
