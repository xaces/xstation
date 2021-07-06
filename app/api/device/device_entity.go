package device

import (
	"xstation/app/api/page"
	"xstation/app/mnger"
	"xstation/models"

	"github.com/wlgd/xutils/orm"
)

// devicePage 分页
type devicePage struct {
	page.Page
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// Where 初始化
func (s *devicePage) Where() *orm.DbWhere {
	var where orm.DbWhere
	if s.StartTime != "" {
		where.Append("created_at >= ?", s.StartTime)
	}
	if s.EndTime != "" {
		where.Append("created_at <= ?", s.EndTime)
	}
	return &where
}

// deviceUpdate 更新
type deviceUpdate struct {
	models.XDeviceOpt
	Id uint64 `json:"id"`
}

func deleteDevices(ids []int) error {
	var devs []models.XDevice
	if _, err := orm.DbFindBy(&devs, "id in (?)", ids); err != nil {
		return err
	}
	if err := orm.DbDeletes(&devs); err != nil {
		return err
	}
	for _, dev := range devs {
		mnger.Dev.Delete(dev.DeviceNo)
	}
	return nil
}
