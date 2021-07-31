package device

import (
	"github.com/wlgd/xutils/orm"
)

// statusPage 分页
type statusPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
	Flag      uint8  `form:"flag"`     // 0-实时 1-补传
	Desc      string `form:"desc"`     //
}

// Where 初始化
func (s *statusPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.Append("device_no like ?", s.DeviceNo)
	where.Append("dtu >= ?", s.StartTime)
	where.Append("dtu <= ?", s.EndTime)
	if s.Flag != 0 {
		where.Append("flag = ?", s.Flag)
	}
	where.Orders = append(where.Orders, s.Desc+" desc")
	return &where
}

// statusGet 获取
type statusGet struct {
	DeviceNo string `form:"deviceNo"` //
	StatusId int64  `form:"statusId"` //
}
