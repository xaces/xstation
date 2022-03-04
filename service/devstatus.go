package service

import (
	"github.com/wlgd/xutils/orm"
)

// StatusPage 分页
type StatusPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
	Flag      int    `form:"flag"`     // 0-实时 1-补传
	Desc      string `form:"desc"`     //
}

// Where 初始化
func (s *StatusPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.String("device_no like ?", s.DeviceNo)
	where.String("dtu >= ?", s.StartTime)
	where.String("dtu <= ?", s.EndTime)
	where.Int("flag = ?", s.Flag)
	where.Orders = append(where.Orders, "dtu desc")
	return &where
}

// StatusGet 获取
type StatusGet struct {
	DeviceNo string `form:"deviceNo"` //
	StatusId uint64 `form:"statusId"` //
}
