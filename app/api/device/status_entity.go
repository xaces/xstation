package device

import (
	"xstation/app/api/page"

	"github.com/wlgd/xutils/orm"
)

// statusPage 分页
type statusPage struct {
	page.Page
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceId  uint64 `form:"deviceId"` //
	Descs     string `form:"descs"`    //
}

// Where 初始化
func (s *statusPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.Append("device_id = ?", s.DeviceId)
	where.Append("dtu >= ?", s.StartTime)
	where.Append("dtu <= ?", s.EndTime)
	where.Orders = append(where.Orders, s.Descs+" desc")
	return &where
}

// statusGet 获取
type statusGet struct {
	DeviceId uint64 `form:"deviceId"` //
	StatusId uint64 `form:"statusId"` //
}