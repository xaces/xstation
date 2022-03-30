package service

import (
	"xstation/model"

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

func DeviceUpdate(m *model.Device, online bool, version, dtype string) error {
	return orm.DbUpdates(m, []string{"version", "type", "online"})
}

func DevcieFindByUser() {

}
