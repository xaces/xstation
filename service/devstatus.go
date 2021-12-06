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
	Flag      uint8  `form:"flag"`     // 0-实时 1-补传
	Desc      string `form:"desc"`     //
}

// Where 初始化
func (s *StatusPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.Append("device_no like ?", s.DeviceNo)
	where.Append("dtu >= ?", s.StartTime+" 00:00:00")
	where.Append("dtu <= ?", s.EndTime+" 23:59:59")
	if s.Flag != 0 {
		where.Append("flag = ?", s.Flag)
	}
	where.Orders = append(where.Orders, "dtu desc")
	return &where
}

// StatusGet 获取
type StatusGet struct {
	DeviceNo string `form:"deviceNo"` //
	StatusId uint64 `form:"statusId"` //
}
