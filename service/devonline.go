package service

import (
	"github.com/wlgd/xutils/orm"
)

// statusPage 分页
type OnlinePage struct {
	orm.DbPage
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
}

// Where 初始化
func (s *OnlinePage) Where() *orm.DbWhere {
	where := s.DbWhere()
	where.String("device_no like ?", s.DeviceNo)
	where.String("on_time >= ?", s.StartTime)
	where.String("off_time <= ?", s.EndTime)
	where.Orders = append(where.Orders, "on_time desc")
	return where
}
