package service

import (
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// DevicePage 分页
type DevicePage struct {
	PageNum    int    `form:"pageNum"`  // 当前页码
	PageSize   int    `form:"pageSize"` // 每页数
	StartTime  string `form:"startTime"`
	EndTime    string `form:"endTime"`
	DeviceNo   string `form:"deviceNo"`
	DeviceName string `form:"deviceName"`
}

// Where 初始化
func (s *DevicePage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.String("created_at >= ?", s.StartTime)
	where.String("created_at <= ?", s.EndTime)
	where.String("no like ?", s.DeviceNo)
	where.String("name like ?", s.DeviceName)
	return &where
}

func DeviceUpdate(m *model.Device, online bool, version, dtype string) error {
	m.Online = online
	if online {
		m.Version = version
		m.Type = dtype
		orm.DbUpdates(m, []string{"version", "type", "last_time", "online"})
	} else {
		orm.DbUpdates(m, []string{"last_time", "online", "last_status"})
	}
	return nil
}
