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
		&model.Device{},
		&model.DevOnline{},
		&model.DevAlarm{},
		&model.DevAlarmDetails{},
		&model.DevAlarmFile{},
		&model.DevCapture{},
		// &model.DevStatus{}, // 本机测试
	)
	orm.CreateTables(&model.DevStatus{})
	task.Timer.AddDbPartFunc(func() {
		orm.NewPartiton(model.DevStatus{}.TableName()).AlterRange("dtu", 30)
	})
	task.Timer.AddDbPartFunc(func() {
		orm.NewPartiton(model.DevAlarmDetails{}.TableName()).AlterRange("dtu", 30)
	})
	task.Timer.AddDbPartFunc(func() {
		orm.NewPartiton(model.DevAlarm{}.TableName()).AlterRange("start_time", 30)
	})
	orm.SetDB(db.Debug())
	return nil
}
