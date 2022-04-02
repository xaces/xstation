package service

import (
	"xstation/model"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

// DevicePage 分页
type DevicePage struct {
	orm.DbPage
	DeviceNo     string `form:"deviceNo"`
	DeviceName   string `form:"deviceName"`
	OrganizeGuid string `form:"organizeGuid"`
	OrganizeId   *int   `form:"organizeId"` // 每页数
}

// Where 初始化
func (s *DevicePage) Where() *orm.DbWhere {
	where := s.DbWhere()
	where.String("device_no like ?", s.DeviceNo)
	where.String("device_name like ?", s.DeviceName)
	where.String("organize_guid = ?", s.OrganizeGuid)
	if s.OrganizeId != nil {
		where.Append("organize_id = ?", *s.OrganizeId)
	}
	return where
}

func DeviceUpdate(m *model.Device, a *xproto.Access) error {
	m.Online = a.Online
	m.Version = a.Version
	m.Type = a.DevType
	m.LastOnlineTime = a.DeviceTime
	if a.Online {
		return orm.DbUpdates(m, "version", "type", "online", "last_online_time")
	}
	return orm.DbUpdates(m, "last_online_time", "online")
}
