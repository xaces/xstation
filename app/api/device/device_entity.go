package device

import (
	"xstation/app/mnger"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// devicePage 分页
type devicePage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
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
	model.DeviceOpt
	Id uint64 `json:"id"`
}

func deleteDevices(ids []int) error {
	var devs []model.Device
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
