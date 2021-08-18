package device

import "github.com/wlgd/xutils/orm"

// statusPage 分页
type onlinePage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
}

// Where 初始化
func (s *onlinePage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.Append("device_no like ?", s.DeviceNo)
	where.Append("on_time >= ?", s.StartTime)
	where.Append("off_time <= ?", s.EndTime)
	where.Orders = append(where.Orders, "on_time desc")
	return &where
}
