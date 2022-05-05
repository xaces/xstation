package db

import (
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

type Option struct {
	Name    string
	Address string
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
	NewPartition("t_devalarm", "start_time", 5).exec()
	return nil
}
