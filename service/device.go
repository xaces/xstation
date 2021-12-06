package service

import (
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// DevicePage 分页
type DevicePage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	No        string `form:"no"`
	Name      string `form:"name"`
}

// Where 初始化
func (s *DevicePage) Where() *orm.DbWhere {
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

// DeviceUpdate 更新
type DeviceUpdate struct {
	model.DeviceOpt
	Id uint64 `json:"id"`
}