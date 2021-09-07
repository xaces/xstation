package device

import (
	"xstation/mnger"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// devicePage 分页
type devicePage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	No        string `form:"no"`
	Name      string `form:"name"`
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
	if s.No != "" {
		where.Append("no like ?", s.No)
	}
	if s.Name != "" {
		where.Append("name like ?", s.Name)
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
		mnger.Dev.Delete(dev.No)
	}
	return nil
}
