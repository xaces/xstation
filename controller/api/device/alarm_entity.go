package device

import "github.com/wlgd/xutils/orm"

// alarmPage 分页
type alarmPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
}

// Where 初始化
func (s *alarmPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.Append("device_no like ?", s.DeviceNo)
	where.Append("start_time >= ?", s.StartTime)
	where.Append("end_time <= ?", s.EndTime)
	where.Orders = append(where.Orders, "start_time desc")
	return &where
}
