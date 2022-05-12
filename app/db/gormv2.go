package db

import (
	"xstation/entity/task"
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
	task.Timer.AddDbPartFunc(func() {
		PartTable(model.DevAlarmDetails{}.TableName()).AlterRange("dtu", 5)
	})
	task.Timer.AddDbPartFunc(func() {
		PartTable(model.DevAlarm{}.TableName()).AlterRange("start_time", 5)
	})
	return nil
}
