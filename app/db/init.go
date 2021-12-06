package db

import (
	"fmt"
	"xstation/app/mnger"
	"xstation/configs"
	"xstation/model"
	"xstation/service"

	"github.com/wlgd/xutils/orm"
)

// loadData 本地数据
func loadData() error {
	var devices []model.Device
	orm.DbFind(&devices)
	if len(devices) > configs.Local.MaxDevNumber {
		return fmt.Errorf("please do not modify the database")
	}
	mnger.Devs.Set(devices)
	return nil
}

func Init() error {
	if err := service.Init(); err != nil {
		return err
	}
	if err := loadData(); err != nil {
		return err
	}
	return nil
}
