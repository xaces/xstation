package mnger

import (
	"fmt"
	"xstation/configs"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

func Init() error {
	// loadData 本地数据
	var devices []model.Device
	orm.DbFind(&devices)
	if len(devices) > configs.Local.MaxDevNumber {
		return fmt.Errorf("please do not modify the database")
	}
	Devs.Set(devices)
	return nil
}
