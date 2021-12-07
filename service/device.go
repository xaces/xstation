package service

import (
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// DevicePage 分页
type DevicePage struct {
	PageNum    int    `form:"pageNum"`  // 当前页码
	PageSize   int    `form:"pageSize"` // 每页数
	StartTime  string `form:"startTime"`
	EndTime    string `form:"endTime"`
	DeviceNo   string `form:"deviceNo"`
	DeviceName string `form:"deviceName"`
}

// Where 初始化
func (s *DevicePage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.String("created_at >= ?", s.StartTime)
	where.String("created_at <= ?", s.EndTime)
	where.String("no like ?", s.DeviceNo)
	where.String("name like ?", s.DeviceName)
	return &where
}

// DeviceUpdate 更新
type DeviceUpdate struct {
	model.DeviceOpt
	Id uint64 `json:"id"`
}
