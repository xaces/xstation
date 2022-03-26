package service

import (
	"github.com/wlgd/xutils/orm"
)

// StatusPage 分页
type StatusPage struct {
	orm.DbPage
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
	Flag      int    `form:"flag"`     // 0-实时 1-补传
	Desc      string `form:"desc"`     //
}

// Where 初始化
func (s *StatusPage) Where() *orm.DbWhere {
	where := s.DbWhere()
	where.String("device_no like ?", s.DeviceNo)
	where.String("dtu >= ?", s.StartTime)
	where.String("dtu <= ?", s.EndTime)
	where.Int("flag = ?", s.Flag)
	where.Orders = append(where.Orders, "dtu desc")
	return where
}