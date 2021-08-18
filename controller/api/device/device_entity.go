package device

import (
	"xstation/mnger"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// devicePage 分页
type devicePage struct {
	PageNum    int    `form:"pageNum"`  // 当前页码
	PageSize   int    `form:"pageSize"` // 每页数
	StartTime  string `form:"startTime"`
	EndTime    string `form:"endTime"`
	DeviceNo   string `form:"deviceNo"`
	DeviceName string `form:"deviceName"`
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
	if s.DeviceNo != "" {
		where.Append("device_no like ?", s.DeviceNo)
	}
	if s.DeviceName != "" {
		where.Append("device_name like ?", s.DeviceName)
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
