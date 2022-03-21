package service

import (
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// statusPage 分页
type OnlinePage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
}

// Where 初始化
func (s *OnlinePage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.String("device_no like ?", s.DeviceNo)
	where.String("on_time >= ?", s.StartTime)
	where.String("off_time <= ?", s.EndTime)
	where.Orders = append(where.Orders, "on_time desc")
	return &where
}

// DevOnlineUpdate 更新链路信息
func DevOnlineUpdate(m *model.DevOnline) error {
	if m.OfflineTime == "" {
		return orm.DbCreate(m)
	}
	return orm.DbUpdateModelBy(m, "guid = ?", m.Guid)
}